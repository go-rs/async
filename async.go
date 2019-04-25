/*!
 * go-rs/async
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package async

type Task func() interface{}
type Channel chan interface{}
type Result []interface{}
type Out interface{}

type Promise struct {
	//TODO:
	//Context - Cancel, Timeout
	//Concurrency vs Parallelism
	//What to do on error? Is it best way to abort all operation?
	//Can I use map for tasks? It will be easy to associate results with keys
}

func (p *Promise) All(tasks []Task) Result {
	workers := make([]Channel, 0, len(tasks))
	for _, task := range tasks {
		result := execute(task)
		workers = append(workers, result)
	}

	// gather data from all channels
	out := make(Result, 0, len(tasks))
	for _, a := range workers {
		out = append(out, <-a)
	}
	return out
}

func execute(task Task) Channel {
	out := make(Channel)
	{
		go func() {
			defer close(out)
			out <- task()
		}()
	}
	return out
}
