package main

import (
	"fmt"
	"go-swim/workerpool"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func main() {
	wp := workerpool.NewSimplePool()

	go func() {
		defer wp.Close()

		squareJob := func(i int) workerpool.Job {
			var args = []any{i}
			return workerpool.Job{Args: args}
		}
		for i := 0; i < 10; i++ {
			wp.Add(squareJob(i))
		}
	}()

	squareWorker := func(j workerpool.Job) workerpool.JobResult {
		if len(j.Args) != 1 {
			return workerpool.JobResult{Error: fmt.Errorf("wrong number of arguments passed")}
		}

		x, ok := j.Args[0].(int)
		if !ok {
			return workerpool.JobResult{Error: fmt.Errorf("expected to get an int, but got: %T", j.Args[0])}
		}

		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

		return workerpool.JobResult{Result: x * x}
	}

	for res := range wp.Run(squareWorker) {
		fmt.Println(res)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		wp.Stop()
	}()
}
