package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func main() {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
	}

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartImmediately()),
		gocron.NewTask(
			func(a string, b int) {
				fmt.Printf("I'm a task %s/%d.", a, b)
			},
			"hello",
			1,
		),
	)
	if err != nil {
		// handle error
	}
	// each job has a unique id
	fmt.Println(j.ID())

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	time.Sleep(time.Minute)

	// when you're done, shut it down
	err = s.Shutdown()
	// or for context-aware teardown:
	// err = s.ShutdownWithContext(ctx)
	if err != nil {
		// handle error
	}
}
