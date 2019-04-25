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
	// buffer channels for go routines
	workers := make(Channel, len(tasks))
	defer close(workers)
	for _, task := range tasks {
		go func(task Task) {
			workers <- task()
		}(task)
	}

	// gather data from all channels
	out := make(Result, 0, len(tasks))
	for i := 0; i < len(tasks); i++ {
		out = append(out, <-workers)
	}
	return out
}
