package taskqueue

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestCompletion(t *testing.T) {
	q := NewQueue(10, 10000, nil)

	s := Task(func(out TaskOutput) {
		time.Sleep(1 * time.Second)
		fmt.Println("done")
	})

	for i := 0; i < 50; i++ {
		q.Enqueue(s)
	}

	q.WaitForCompletion()
}

func TestSequenceTaskOutput(t *testing.T) {
	out := NewSequenceTaskOutput(10000)
	q := NewQueue(10, 10000, out)

	s := Task(func(out TaskOutput) {
		out.Write(rand.Int31())
	})

	for i := 0; i < 1000; i++ {
		q.Enqueue(s)
	}

	go q.WaitForCompletion()

	out.WaitForOutput(TaskOutputHandler(func(v interface{}) {
		m, _ := v.(int32)
		fmt.Println(m)
	}))
}
