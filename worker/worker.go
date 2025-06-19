package worker

import (
	"context"
	"fmt"
	"log"
	"time"
)

func StartWorker(ctx context.Context, id int, jobs <-chan Job, results chan<- string) {
	logger := log.New(log.Writer(), fmt.Sprintf("[Worker %d] ", id), log.LstdFlags)

	for {
		select {
		case <-ctx.Done():
			logger.Println("Context cancelled. Shutting down.")
			return
		case job, ok := <-jobs:
			if !ok {
				logger.Println("Job channel closed. Exiting.")
				return
			}

			logger.Printf("Started job %d with payload: %s\n", job.ID, job.Payload)

			// Simulate work with cancellation check
			select {
			case <-ctx.Done():
				logger.Printf("Job %d cancelled before completion.\n", job.ID)
				return
			case <-time.After(1 * time.Second):
				// Simulated work done
			}

			result := fmt.Sprintf("Finished job %d with payload: %s", job.ID, job.Payload)
			logger.Println(result)
			results <- result
		}
	}
}
