package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		out := make(Bi)
		close(out)
		return out
	}

	if done == nil {
		pipeline := stages[0](in)

		for i := 1; i < len(stages); i++ {
			pipeline = stages[i](pipeline)
		}

		return pipeline
	}

	chanAdapter := func(in In, done In) Out {
		out := make(Bi)
		go func() {
			stopAndDrain := func() {
				close(out)
				for v := range in {
					_ = v
				}
			}

			for {
				select {
				case <-done:
					stopAndDrain()
					return
				case v, ok := <-in:
					if !ok {
						close(out)
						return
					}
					select {
					case <-done:
						stopAndDrain()
						return
					case out <- v:
					}
				}
			}
		}()
		return out
	}

	pipeline := stages[0](chanAdapter(in, done))

	for i := 1; i < len(stages); i++ {
		pipeline = stages[i](chanAdapter(pipeline, done))
	}

	return chanAdapter(pipeline, done)
}
