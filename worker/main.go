package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ktbsomen/gobullmq"
	"github.com/ktbsomen/gobullmq/types"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	params, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	workerClient := redis.NewClient(params)

	workerProcess := func(ctx context.Context, job *types.Job, api gobullmq.WorkerProcessAPI) (interface{}, error) {
		fmt.Printf("Processing job: %s\n", job.Name)
		return "ok", nil
	}

	worker, err := gobullmq.NewWorker(ctx, "process-asset", gobullmq.WorkerOptions{
		Concurrency:     1,
		StalledInterval: 300000000,
		Backoff:         &gobullmq.BackoffOptions{Type: "exponential", Delay: 500},
	}, workerClient, workerProcess)
	if err != nil {
		log.Fatal(err)
	}

	// Run blocks until ctx is cancelled
	if err := worker.Run(); err != nil {
		log.Printf("Worker error: %v", err)
	}
}
