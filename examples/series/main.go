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
			fmt.Println("Running Hello..........")
			time.Sleep(2 * time.Second)
			fmt.Println("Completed Hello..........")
			//return errors.New("innocent")
			return "Hello"
		},
		func() interface{} {
			fmt.Println("Running World..........")
			time.Sleep(1 * time.Second)
			fmt.Println("Completed World..........")
			return "World"
		},
	}

	fmt.Println("From: ", time.Now())
	result, err := promise.Series(tasks)
	fmt.Println("End: ", time.Now())

	fmt.Println("Result: ", result, err)

	fmt.Println("exit")
}
