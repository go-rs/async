package main

import (
	"fmt"
	"time"

	"github.com/go-rs/async"
)

func main() {
	fmt.Println("start")
	var promise async.Promise

	tasks := async.Tasks{
		func() interface{} {
			fmt.Println("Running Hello..........")
			time.Sleep(4 * time.Second)
			fmt.Println("Completed Hello..........")
			return "Hello"
		},
		func() interface{} {
			fmt.Println("Running World..........")
			time.Sleep(2 * time.Second)
			fmt.Println("Completed World..........")
			return "World"
		},
	}

	fmt.Println("From: ", time.Now())
	result := promise.Parallel(tasks)
	fmt.Println("End: ", time.Now())

	fmt.Println("Result: ", result)

	fmt.Println("=============================================================")

	mtasks := async.ETasks{
		"Hello": func() (interface{}, error) {
			fmt.Println("Running Hello..........")
			time.Sleep(4 * time.Second)
			fmt.Println("Completed Hello..........")
			return "Hello", nil //errors.New("EHello")
		},
		"World": func() (interface{}, error) {
			fmt.Println("Running World..........")
			time.Sleep(2 * time.Second)
			fmt.Println("Completed World..........")
			return "World", nil //errors.New("EWorld")
		},
	}

	fmt.Println("From: ", time.Now())
	mresult, err := promise.ParallelWithMap(mtasks)
	fmt.Println("End: ", time.Now())

	fmt.Println("Result: ", mresult, err)

	fmt.Println("exit")
}
