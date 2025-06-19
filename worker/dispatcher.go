package worker

import "fmt"

func DispatchJobs(jobCount int) <-chan Job {
	jobs := make(chan Job, jobCount)
	go func() {
		for i := 1; i <= jobCount; i++ {
			jobs <- Job{ID: i, Payload: fmt.Sprintf("Payload %d", i)}
		}
		close(jobs)
	}()
	return jobs
}
