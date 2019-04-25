/*!
 * go-rs/async
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package async

type Task func() interface{}
type Result chan interface{}

type Promise struct {
	//TODO:
	//Context - Cancel, Timeout
	//Concurrency vs Parallelism
	//What to do on error? Is it best way to abort all operation?
}

func (p *Promise) All(tasks []Task) []Result {
	workers := make([]Result, 0, len(tasks))
	for _, task := range tasks {
		result := execute(task)
		workers = append(workers, result)
	}
	return workers
}

func execute(task Task) Result {
	out := make(Result)
	{
		go func() {
			defer close(out)
			out <- task()
		}()
	}
	return out
}
