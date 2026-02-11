package main

import (
	"fmt"
	"strings"
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

How it works (step by step):
─────────────────────────────
1. You create a JOBS channel (the queue) and a RESULTS channel (to collect output).
2. You spin up N worker goroutines. Each one loops over the jobs channel:
       for job := range jobs { ... }
   When the channel is closed and drained, the loop exits automatically.
3. The main goroutine pushes jobs into the jobs channel, then closes it.
4. A separate goroutine waits for all workers to finish (via sync.WaitGroup),
   then closes the results channel so the collector loop can exit.
5. The main goroutine (or another consumer) ranges over results to collect output.

Diagram:
────────
                    ┌──────────┐
           ┌──────▶│ Worker 1  │──────┐
           │       └──────────┘      │
  ┌──────┐ │       ┌──────────┐      │  ┌─────────┐
  │ Jobs │─┼──────▶│ Worker 2  │──────┼─▶│ Results │
  │ Chan │ │       └──────────┘      │  │  Chan   │
  └──────┘ │       ┌──────────┐      │  └─────────┘
           └──────▶│ Worker 3  │──────┘
                    └──────────┘

  - All workers READ from the same jobs channel (Go distributes automatically).
  - All workers WRITE to the same results channel.
  - WaitGroup tracks when every worker is done.

Key Components:
- Job Queue: Channel that holds tasks to be processed
- Workers: Fixed number of goroutines that consume jobs from the queue
- Results Channel: Optional channel to collect results

Why buffered channels?
──────────────────────
  jobs    := make(chan Job, numJobs)      // buffered so sender won't block
  results := make(chan Result, numJobs)   // buffered so workers won't block

  A buffered channel lets the sender push items without waiting for a receiver,
  up to the buffer size. This decouples producers from consumers.

Benefits:
- Limits concurrent operations (prevents resource exhaustion)
- Reuses goroutines (reduces overhead of creating/destroying goroutines)
- Better resource management and predictable performance
- Back-pressure: if all workers are busy, the jobs channel blocks the sender

Use Cases:
- Processing large datasets in parallel
- Handling multiple API requests with rate limiting
- Image/video processing pipelines
- Database batch operations
- Web scraping with controlled concurrency

Best Practices:
1. Choose worker count based on CPU cores or I/O constraints
2. Use buffered channels for job queue to prevent blocking
3. Always close job channel when done sending jobs
4. Use sync.WaitGroup to wait for all workers to finish
5. Handle panics in worker goroutines to prevent pool shutdown
*/

// ─────────────────────────────────────────────────────────────
// TYPES
// ─────────────────────────────────────────────────────────────

// Job represents a unit of work that a worker will process.
type Job struct {
	ID   int
	Data string
}

// Result represents the output produced by a worker for a given Job.
type Result struct {
	JobID  int
	Output string
	Error  error
}

// ─────────────────────────────────────────────────────────────
// WORKER  (runs as a goroutine)
// ─────────────────────────────────────────────────────────────

// poolWorker reads from the shared jobs channel, processes each job, and
// sends the result to the results channel. It calls wg.Done() when the
// jobs channel is closed and drained (loop exits).
func poolWorker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d started  job %d\n", id, job.ID)

		// Simulate work
		time.Sleep(500 * time.Millisecond)

		// Process the job
		result := Result{
			JobID:  job.ID,
			Output: fmt.Sprintf("Processed by worker %d: %s", id, job.Data),
			Error:  nil,
		}

		results <- result
		fmt.Printf("Worker %d finished job %d\n", id, job.ID)
	}
}

// ─────────────────────────────────────────────────────────────
// DEMO
// ─────────────────────────────────────────────────────────────

// DemoWorkerPools demonstrates the worker pool pattern.
func DemoWorkerPools() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("WORKER POOLS DEMONSTRATION")
	fmt.Println(strings.Repeat("=", 60) + "\n")

	const (
		numWorkers = 3
		numJobs    = 10
	)

	// ── Step 1: Create channels ──────────────────────────────
	jobs := make(chan Job, numJobs)       // buffered job queue
	results := make(chan Result, numJobs) // buffered results collector

	// ── Step 2: Create WaitGroup for workers ─────────────────
	var wg sync.WaitGroup

	// ── Step 3: Start workers ────────────────────────────────
	fmt.Printf("Starting %d workers...\n\n", numWorkers)
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go poolWorker(w, jobs, results, &wg)
	}

	// ── Step 4: Send jobs ────────────────────────────────────
	fmt.Printf("Sending %d jobs to the pool...\n\n", numJobs)
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{
			ID:   j,
			Data: fmt.Sprintf("Job-%d-Data", j),
		}
	}
	close(jobs) // Close channel to signal no more jobs

	// ── Step 5: Wait for workers & close results ─────────────
	go func() {
		wg.Wait()
		close(results)
	}()

	// ── Step 6: Collect results ──────────────────────────────
	fmt.Println("\nCollecting results:\n")
	for result := range results {
		if result.Error != nil {
			fmt.Printf("Job %d failed: %v\n", result.JobID, result.Error)
		} else {
			fmt.Printf("Job %d result: %s\n", result.JobID, result.Output)
		}
	}

	fmt.Println("\n✅ All jobs processed successfully!")

	fmt.Println("\nKey Takeaways:")
	fmt.Println("- Workers process jobs concurrently")
	fmt.Println("- Limited number of workers (resource control)")
	fmt.Println("- Jobs distributed automatically among available workers")
	fmt.Println("- Results collected in order they complete (not job order)")
}
