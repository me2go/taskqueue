package taskqueue

import (
	"sync"
)

type Task func(TaskOutput)

type TaskOutput interface {
	Write(interface{})
	Close()
}

type TaskQueue struct {
	sync.WaitGroup

	worker chan struct{}
	pool   chan Task
	close  chan struct{}
	output TaskOutput
}

func NewQueue(workerNum int, queueSize int, output TaskOutput) *TaskQueue {
	q := &TaskQueue{
		worker: make(chan struct{}, workerNum),
		pool:   make(chan Task, queueSize),
		close:  make(chan struct{}),
		output: &NullTaskOutput{},
	}

	if output != nil {
		q.output = output
	}

	//init worker pool
	for i := 0; i < workerNum; i++ {
		q.worker <- struct{}{}
	}

	q.Add(1)
	go q.run()

	return q
}

func (q *TaskQueue) Enqueue(t Task) {
	q.pool <- t
}

func (q *TaskQueue) run() {
	defer q.Done()
	for {
		select {
		case t := <-q.pool:
			if t == nil {
				continue
			}
			q.exec(t)

		case <-q.close:
			//queue has closed
			return
		}
	}
}

func (q *TaskQueue) exec(t Task) {
	//choose a worker
	<-q.worker
	//run task
	q.Add(1)
	go func(tt Task) {
		defer q.Done()
		tt(q.output)
		q.worker <- struct{}{}
	}(t)
}

func (q *TaskQueue) WaitForCompletion() {
	//shutdown task pool
	close(q.pool)
	//shutdown daemon goroutine
	close(q.close)
	//drain the task pool
	for t := range q.pool {
		q.exec(t)
	}
	//wait all goroutine exit
	q.Wait()
	close(q.worker)
	q.output.Close()
}
