/*!
 * go-rs/async
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package async

type Tasks []func() interface{}
type ETasks map[string]func() (interface{}, error)

// Reference: https://blog.golang.org/pipelines
// Reference: https://en.wikipedia.org/wiki/Moore%27s_law
// Reference: https://en.wikipedia.org/wiki/Amdahl%27s_law

type Promise struct {
	//TODO: Context - Cancel, Timeout
	//Rob Pike: Concurrency enables Parallelism
	//What to do on error? Is it best way to abort all operation?
}

//
// parallel execution of all tasks
//
func (p *Promise) Parallel(tasks []func() interface{}) []interface{} {
	// can not use buffer channels, because not able to maintain output sequence
	// using slice channels
	workers := make([]chan interface{}, 0, len(tasks))
	for _, task := range tasks {
		workers = append(workers, func() chan interface{} {
			out := make(chan interface{})
			go func(task func() interface{}) {
				defer close(out)
				out <- task()
			}(task)
			return out
		}())
	}

	// gather data from all channels
	out := make([]interface{}, 0, len(tasks))
	for _, result := range workers {
		out = append(out, <-result)
	}

	return out
}

//
// Parallel execution of all tasks, but returns result as well as error
//
func (p *Promise) ParallelWithMap(tasks map[string]func() (interface{}, error)) (map[string]interface{}, map[string]error) {
	// can not use buffer channels, because not able to maintain output sequence
	// using slice channels
	workers := make(map[string]chan interface{})
	for key, task := range tasks {
		workers[key] = func() chan interface{} {
			out := make(chan interface{})
			go func(task func() (interface{}, error)) {
				defer close(out)
				val, err := task()
				if err != nil {
					out <- err
				} else {
					out <- val
				}
			}(task)
			return out
		}()
	}

	// gather data from all channels
	out := make(map[string]interface{})
	err := make(map[string]error)
	for key, result := range workers {
		var _err error
		_val := <-result
		if _err = isError(_val); _err != nil {
			err[key] = _err
			_val = nil
		}
		out[key] = _val
	}

	flag := true
	for _, e := range err {
		if e != nil {
			flag = false
			break
		}
	}
	if flag {
		err = nil
	}

	return out, err
}

func isError(val interface{}) (err error) {
	switch val.(type) {
	case error:
		err = val.(error)
	}
	return
}
