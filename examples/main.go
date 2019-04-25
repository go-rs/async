package main

import (
	"fmt"
	"time"

	"github.com/go-rs/async"
)

func main() {
	fmt.Println("start")
	var promise async.Promise

	tasks := []async.Task{
		func() interface{} {
			time.Sleep(1000 * time.Microsecond)
			return "Hello"
		},
		func() interface{} {
			time.Sleep(100 * time.Microsecond)
			return "World"
		},
	}

	for _, a := range promise.All(tasks) {
		// switch statement for results
		fmt.Println(<-a)
	}

	fmt.Println("exit")
}
