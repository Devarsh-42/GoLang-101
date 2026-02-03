package main

import (
	"fmt"
	"sync"
	"time"
)

/*
WORKER POOLS
============

What is it?
-----------
A Worker Pool is a concurrency pattern where a fixed number of goroutines (workers)
process jobs from a shared channel. This pattern helps control the number of concurrent
operations and prevents resource exhaustion.

Key Components:
- Job Queue: Channel that holds tasks to be processed
- Workers: Fixed number of goroutines that consume jobs from the queue
- Results Channel: Optional channel to collect results

Benefits:
- Limits concurrent operations (prevents resource exhaustion)
- Reuses goroutines (reduces overhead of creating/destroying goroutines)
- Better resource management and predictable performance

Use Cases:
- Processing large datasets in parallel
- Handling multiple API requests with rate limiting
- Image/video processing pipelines
- Database batch operations
- Web scraping with controlled concurrency






































































































}	fmt.Println("- Results collected in order they complete (not job order)")	fmt.Println("- Jobs distributed automatically among available workers")	fmt.Println("- Limited number of workers (resource control)")	fmt.Println("- Workers process jobs concurrently")	fmt.Println("\nKey Takeaways:")	fmt.Println("\n✅ All jobs processed successfully!")	}		}			fmt.Printf("Job %d result: %s\n", result.JobID, result.Output)		} else {			fmt.Printf("Job %d failed: %v\n", result.JobID, result.Error)		if result.Error != nil {	for result := range results {	fmt.Println("\nCollecting results:\n")	// Collect results	}()		close(results)		wg.Wait()	go func() {	// Wait for all workers to finish and close results	close(jobs) // Close channel to signal no more jobs	}		}			Data: fmt.Sprintf("Job-%d-Data", j),			ID:   j,		jobs <- Job{	for j := 1; j <= numJobs; j++ {	fmt.Printf("Sending %d jobs to the pool...\n\n", numJobs)	// Send jobs	}		go worker(w, jobs, results, &wg)		wg.Add(1)	for w := 1; w <= numWorkers; w++ {	fmt.Printf("Starting %d workers...\n\n", numWorkers)	// Start workers	var wg sync.WaitGroup	// Create WaitGroup for workers	results := make(chan Result, numJobs)	jobs := make(chan Job, numJobs)	// Create channels	)		numJobs    = 10		numWorkers = 3	const (	fmt.Println("="*60 + "\n")	fmt.Println("WORKER POOLS DEMONSTRATION")	fmt.Println("\n" + "="*60)func DemoWorkerPools() {// DemoWorkerPools demonstrates the worker pool pattern}	}		fmt.Printf("Worker %d finished job %d\n", id, job.ID)		results <- result				}			Error:  nil,			Output: fmt.Sprintf("Processed by worker %d: %s", id, job.Data),			JobID:  job.ID,		result := Result{		// Process the job				time.Sleep(500 * time.Millisecond)		// Simulate work				fmt.Printf("Worker %d started job %d\n", id, job.ID)	for job := range jobs {		defer wg.Done()func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {// Worker processes jobs from the jobs channel}	Error  error	Output string	JobID  inttype Result struct {// Result represents the output of a job}	Data string	ID   inttype Job struct {// Job represents a unit of work*/5. Handle panics in worker goroutines to prevent pool shutdown4. Use sync.WaitGroup to wait for all workers to finish3. Always close job channel when done sending jobs2. Use buffered channels for job queue to prevent blocking1. Choose worker count based on CPU cores or I/O constraintsBest Practices: