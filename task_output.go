package taskqueue

type TaskOutputHandler func(interface{})

type NullTaskOutput struct{}

func (ntq *NullTaskOutput) Write(out interface{}) {}
func (ntq *NullTaskOutput) Close()                {}

type SequenceTaskOutput struct {
	outputPool chan interface{}
}

func NewSequenceTaskOutput(size int) *SequenceTaskOutput {
	return &SequenceTaskOutput{
		outputPool: make(chan interface{}, size),
	}
}

func (st *SequenceTaskOutput) Write(out interface{}) {
	st.outputPool <- out
}

func (st *SequenceTaskOutput) Close() {
	close(st.outputPool)
}

func (st *SequenceTaskOutput) WaitForOutput(h TaskOutputHandler) {
	for out := range st.outputPool {
		h(out)
	}
}
