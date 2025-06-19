package helper

import (
	"encoding/json"
	"fmt"
	"go-jwt-api/worker"
	"net/http"
)

func TriggerWorker(w http.ResponseWriter, r *http.Request) {
	// Default values
	const defNumJobs = 5
	const defNumWorkers = 3

	var req struct {
		JobCount    int `json:"job_count"`
		WorkerCount int `json:"worker_count"`
	}

	numJobs := defNumJobs
	numWorkers := defNumWorkers

	// Try to decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Invalid JSON. Using default values.")
	} else {
		// Use values from request if decoding succeeds
		if req.JobCount > 0 {
			numJobs = req.JobCount
		}
		if req.WorkerCount > 0 {
			numWorkers = req.WorkerCount
		}
	}

	// Dispatch jobs and start workers
	jobs := worker.DispatchJobs(numJobs)
	results := make(chan string, numJobs)

	for w := 1; w <= numWorkers; w++ {
		go worker.StartWorker(r.Context(), w, jobs, results)
	}

	var output []string
	for a := 1; a <= numJobs; a++ {
		output = append(output, <-results)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)

}
