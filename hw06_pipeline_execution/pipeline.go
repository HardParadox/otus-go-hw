package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	for _, stage := range stages {
		if stage == nil {
			continue
		}

		out := make(Bi)

		go func(in In, done In) {
			defer close(out)
			for {
				select {
				case <-done:
					return
				case val, ok := <-in:
					if !ok {
						return
					}
					out <- val
				}
			}
		}(in, done)

		in = stage(out)
	}

	return in
}
