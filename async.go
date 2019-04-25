/*!
 * go-rs/async
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package async

type Task func() interface{}
type Channel chan interface{}
type Result []interface{}

type Promise struct {
	//TODO:
	//Context - Cancel, Timeout
	//Concurrency vs Parallelism
	//What to do on error? Is it best way to abort all operation?
	//Can I use map for tasks? It will be easy to associate results with keys
}

func (p *Promise) All(tasks []Task) Result {
	// can not use buffer channels, because not able to maintain sequence
	// using channel slice
	workers := make([]Channel, 0, len(tasks))
	for _, task := range tasks {
		workers = append(workers, func() Channel {
			out := make(Channel)
			go func(task Task) {
				defer close(out)
				out <- task()
			}(task)
			return out
		}())
	}

	// gather data from all channels
	out := make(Result, 0, len(tasks))
	for _, result := range workers {
		out = append(out, <-result)
	}
	return out
}
