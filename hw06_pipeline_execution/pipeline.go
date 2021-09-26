package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func Done(done In, valueStream In) Out {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for {
			select {
			case <-done:
				return
			case val, ok := <-valueStream:
				if !ok {
					return
				}
				takeStream <- val
			}
		}
	}()
	return takeStream
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	lin := in
	for _, stage := range stages {
		select {
		case <-done:
			return nil
		default:
			lin = Done(done, stage(lin))
		}
	}
	return lin
}
