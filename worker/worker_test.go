package worker_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"go-jwt-api/worker"
)

func TestDispatchJobs(t *testing.T) {
	jobCount := 5
	jobs := worker.DispatchJobs(jobCount)

	receivedJobs := []worker.Job{}
	for job := range jobs {
		receivedJobs = append(receivedJobs, job)
	}

	if len(receivedJobs) != jobCount {
		t.Errorf("expected %d jobs, got %d", jobCount, len(receivedJobs))
	}

	for i, job := range receivedJobs {
		expectedID := i + 1
		if job.ID != expectedID {
			t.Errorf("expected job ID %d, got %d", expectedID, job.ID)
		}
		if !strings.HasPrefix(job.Payload, "Payload") {
			t.Errorf("unexpected payload: %s", job.Payload)
		}
	}
}

func TestStartWorker(t *testing.T) {
	jobCount := 3
	jobs := make(chan worker.Job, jobCount)
	results := make(chan string, jobCount)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker
	go worker.StartWorker(ctx, 1, jobs, results)

	// Send jobs
	for i := 1; i <= jobCount; i++ {
		jobs <- worker.Job{ID: i, Payload: "test payload"}
	}
	close(jobs)

	received := []string{}
	for i := 0; i < jobCount; i++ {
		select {
		case res := <-results:
			received = append(received, res)
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for job result")
		}
	}

	if len(received) != jobCount {
		t.Errorf("expected %d results, got %d", jobCount, len(received))
	}

	for _, res := range received {
		if !strings.Contains(res, "Finished job") {
			t.Errorf("unexpected result: %s", res)
		}
	}
}

func TestStartWorker_ContextCancel(t *testing.T) {
	jobCount := 2
	jobs := make(chan worker.Job, jobCount)
	results := make(chan string, jobCount)
	ctx, cancel := context.WithCancel(context.Background())

	// Start worker
	go worker.StartWorker(ctx, 1, jobs, results)

	// Send one job
	jobs <- worker.Job{ID: 1, Payload: "test payload"}

	// Cancel context before job finishes
	cancel()

	// Close jobs to allow worker to exit
	close(jobs)

	// Give some time for worker to handle cancellation
	time.Sleep(100 * time.Millisecond)

	// Since context cancelled, results may be empty or partial
	// We just check no deadlock and function returns
}
