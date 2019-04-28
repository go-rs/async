/*!
 * go-rs/async
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package async

type Tasks []func() interface{}
type ETasks map[string]func() (interface{}, error)

type channel chan interface{}

// Reference: https://blog.golang.org/pipelines
// Reference: https://en.wikipedia.org/wiki/Moore%27s_law
// Reference: https://en.wikipedia.org/wiki/Amdahl%27s_law
// Rob Pike: Concurrency enables Parallelism

type Promise struct {
	//TODO: Context - Cancel, Timeout
}

//
// parallel execution of all tasks
//
func (p *Promise) Parallel(tasks Tasks) []interface{} {
	// can not use buffer channels, because not able to maintain output sequence
	// using slice channels
	workers := make([]channel, 0, len(tasks))
	for _, task := range tasks {
		workers = append(workers, func() channel {
			out := make(channel)
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
func (p *Promise) ParallelWithMap(tasks ETasks) (map[string]interface{}, map[string]error) {
	// can not use buffer channels, because not able to maintain output sequence
	// using map channels
	workers := make(map[string]channel)
	for key, task := range tasks {
		workers[key] = func() channel {
			out := make(channel)
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

	// set error as nil if nothing caught
	err = hasError(err)

	return out, err
}

// check val type, if error then return cast to error
func isError(val interface{}) (err error) {
	switch val.(type) {
	case error:
		err = val.(error)
	}
	return
}

// check map and look for an error, otherwise set as nil
func hasError(err map[string]error) map[string]error {
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

	return err
}
