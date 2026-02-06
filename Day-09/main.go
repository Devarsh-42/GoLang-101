package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// =============================================================================
// 📚 CONTEXT PACKAGE - All Important Components
// =============================================================================
//
// The `context` package is used to:
//   1. Cancel long-running operations
//   2. Set deadlines/timeouts
//   3. Pass request-scoped values (like user ID, trace ID)
//
// IMPORTANT RULES:
//   - Always pass context as the FIRST parameter of a function
//   - Never store context in a struct
//   - Use context.Background() as the root (top-level) context
//   - Use context.TODO() when you're not sure which context to use yet
//
// =============================================================================

// Job represents a unit of work for the worker pool
type Job struct {
	ID       int
	Duration time.Duration // simulates how long the work takes
}

// Result represents the output of a completed job
type Result struct {
	JobID    int
	WorkerID int
	Output   string
	Duration time.Duration
}

func main() {
	fmt.Println("========================================")
	fmt.Println("  Go Context Package - All Components")
	fmt.Println("========================================\n")

	// 1. context.Background() - the root context, never cancelled
	ctx := context.Background()

	// Run all examples
	exampleWithValue(ctx)
	fmt.Println()
	exampleWithTimeout(ctx)
	fmt.Println()
	exampleWithDeadline(ctx)
	fmt.Println()
	exampleWithCancel(ctx)
	fmt.Println()
	exampleWorkerPool(ctx)
}

// =============================================================================
// 1️⃣ context.WithValue - Pass data through context (like request ID, user ID)
// =============================================================================
// Syntax: ctx := context.WithValue(parentCtx, key, value)
// Use case: Pass request-scoped data (NOT for function parameters!)

type contextKey string // custom type to avoid key collisions

const userIDKey contextKey = "userID"

func exampleWithValue(ctx context.Context) {
	fmt.Println("--- 1. context.WithValue ---")

	// Attach a value to the context
	ctxWithUser := context.WithValue(ctx, userIDKey, "user-123")

	// Pass context to another function
	processRequest(ctxWithUser)
}

func processRequest(ctx context.Context) {
	// Retrieve the value from context
	userID := ctx.Value(userIDKey)
	if userID != nil {
		fmt.Printf("  Processing request for user: %s\n", userID)
	} else {
		fmt.Println("  No user ID found in context")
	}
}

// =============================================================================
// 2️⃣ context.WithTimeout - Auto-cancel after a duration
// =============================================================================
// Syntax: ctx, cancel := context.WithTimeout(parentCtx, duration)
// Use case: HTTP requests, database queries, any operation with a time limit

func exampleWithTimeout(ctx context.Context) {
	fmt.Println("--- 2. context.WithTimeout ---")

	// Create context that auto-cancels after 1 second
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel() // ALWAYS call cancel to release resources

	// Simulate work that takes 2 seconds (will timeout!)
	select {
	case <-time.After(2 * time.Second):
		fmt.Println("  Work finished (won't reach here)")
	case <-ctxWithTimeout.Done():
		fmt.Printf("  Timeout! Error: %v\n", ctxWithTimeout.Err())
	}
}

// =============================================================================
// 3️⃣ context.WithDeadline - Cancel at a specific time
// =============================================================================
// Syntax: ctx, cancel := context.WithDeadline(parentCtx, time.Time)
// Use case: When you need to finish by a specific clock time

func exampleWithDeadline(ctx context.Context) {
	fmt.Println("--- 3. context.WithDeadline ---")

	// Set deadline to 500ms from now
	deadline := time.Now().Add(500 * time.Millisecond)
	ctxWithDeadline, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("  Work finished (won't reach here)")
	case <-ctxWithDeadline.Done():
		fmt.Printf("  Deadline passed! Error: %v\n", ctxWithDeadline.Err())
	}
}

// =============================================================================
// 4️⃣ context.WithCancel - Manually cancel when you want
// =============================================================================
// Syntax: ctx, cancel := context.WithCancel(parentCtx)
// Use case: Stop goroutines when you no longer need their results

func exampleWithCancel(ctx context.Context) {
	fmt.Println("--- 4. context.WithCancel ---")

	ctxWithCancel, cancel := context.WithCancel(ctx)

	// Start a goroutine that runs until cancelled
	go func() {
		for i := 1; ; i++ {
			select {
			case <-ctxWithCancel.Done():
				fmt.Printf("  Goroutine stopped! Ran %d iterations. Error: %v\n", i-1, ctxWithCancel.Err())
				return
			default:
				fmt.Printf("  Working... iteration %d\n", i)
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	// Let it run for 500ms, then cancel
	time.Sleep(500 * time.Millisecond)
	cancel()                           // This stops the goroutine
	time.Sleep(100 * time.Millisecond) // small wait to let goroutine print
}

// =============================================================================
// 5️⃣ WORKER POOL WITH CONTEXT - Real-world pattern!
// =============================================================================
//
// What is a Worker Pool?
//   - A fixed number of goroutines (workers) that process jobs from a shared queue
//   - Instead of spawning 1000 goroutines for 1000 jobs, you use 3-5 workers
//   - Workers pick up jobs from a channel, process them, and send results back
//
// Why use it?
//   - Controls resource usage (CPU, memory, connections)
//   - Prevents system overload
//   - Context lets you cancel ALL workers at once
//
// How it works:
//   ┌──────────┐    jobs channel    ┌──────────┐    results channel    ┌─────────┐
//   │  Main    │ ──────────────────▶│ Workers  │ ────────────────────▶│ Collect │
//   │ (sends   │                    │ (process │                      │ Results │
//   │  jobs)   │                    │  jobs)   │                      │         │
//   └──────────┘                    └──────────┘                      └─────────┘
//         │                              ▲
//         │         context cancel        │
//         └───────── stops all ──────────┘
//

func exampleWorkerPool(ctx context.Context) {
	fmt.Println("--- 5. Worker Pool with Context ---")

	const numWorkers = 3 // number of workers
	const numJobs = 10   // total jobs to process

	// Create a context with timeout - pool must finish in 3 seconds
	ctxPool, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Channels for communication
	jobs := make(chan Job, numJobs)       // buffered channel for jobs
	results := make(chan Result, numJobs) // buffered channel for results

	// WaitGroup to know when all workers are done
	var wg sync.WaitGroup

	// Start workers
	fmt.Printf("  Starting %d workers...\n", numWorkers)
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(ctxPool, w, jobs, results, &wg)
	}

	// Send jobs to the jobs channel
	fmt.Printf("  Sending %d jobs...\n", numJobs)
	for j := 1; j <= numJobs; j++ {
		// Each job takes a random time between 100ms-800ms
		duration := time.Duration(100+rand.Intn(700)) * time.Millisecond
		jobs <- Job{ID: j, Duration: duration}
	}
	close(jobs) // close channel to signal no more jobs

	// Wait for all workers to finish, then close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and print results
	fmt.Println("\n  📊 Results:")
	completed := 0
	for result := range results {
		fmt.Printf("    ✅ Job #%d completed by Worker #%d in %v - %s\n",
			result.JobID, result.WorkerID, result.Duration, result.Output)
		completed++
	}

	fmt.Printf("\n  Completed %d/%d jobs\n", completed, numJobs)
	if ctxPool.Err() != nil {
		fmt.Printf("  ⚠️  Pool stopped early: %v\n", ctxPool.Err())
	}
}

// worker processes jobs from the jobs channel until cancelled or no more jobs
func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			// Context was cancelled (timeout or manual cancel)
			fmt.Printf("  🛑 Worker #%d stopped: %v\n", id, ctx.Err())
			return

		case job, ok := <-jobs:
			if !ok {
				// Jobs channel is closed, no more work
				fmt.Printf("  🏁 Worker #%d finished (no more jobs)\n", id)
				return
			}

			// Simulate doing work (check context during work too!)
			select {
			case <-ctx.Done():
				fmt.Printf("  🛑 Worker #%d cancelled during job #%d\n", id, job.ID)
				return
			case <-time.After(job.Duration):
				// Job completed successfully
				results <- Result{
					JobID:    job.ID,
					WorkerID: id,
					Output:   fmt.Sprintf("Processed in %v", job.Duration),
					Duration: job.Duration,
				}
			}
		}
	}
}
