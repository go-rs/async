/*!
 * go-rs/async
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package async

type Task func() interface{}
type STask func(interface{}) interface{}
type Channel chan interface{}
type Result []interface{}

// Reference: https://blog.golang.org/pipelines
// Reference: https://en.wikipedia.org/wiki/Moore%27s_law
// https://en.wikipedia.org/wiki/Amdahl%27s_law

type Promise struct {
	//TODO:
	//Context - Cancel, Timeout
	//Concurrency vs Parallelism
	//What to do on error? Is it best way to abort all operation?
	//Can I use map for tasks? It will be easy to associate results with keys
}

//
// parallel execution of all tasks
//
func (p *Promise) Parallel(tasks []Task) Result {
	// can not use buffer channels, because not able to maintain output sequence
	// using slice channels
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

//
// sequential execution of all tasks
//
func (p *Promise) Series(tasks []Task) (out interface{}, err error) {
	// unbuffered channel
	worker := make(Channel)
	defer close(worker)

	for _, task := range tasks {
		// need to check the benefits of executing tasks on another go routine
		go func(task Task) {
			worker <- task()
		}(task)

		out = <-worker

		if err = isError(out); err != nil {
			out = nil
			break
		}
	}

	return
}

//
// sequential execution of all tasks, but output of all tasks will be input for next task
//
func (p *Promise) Waterfall(tasks []STask) (out interface{}, err error) {
	// unbuffered channel
	worker := make(Channel)
	defer close(worker)

	for _, task := range tasks {

		go func(task STask) {
			worker <- task(out)
		}(task)

		out = <-worker

		if err = isError(out); err != nil {
			out = nil
			break
		}
	}

	return
}

func isError(val interface{}) (err error) {
	switch val.(type) {
	case error:
		err = val.(error)
	}
	return
}
