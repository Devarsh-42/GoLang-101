package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

/*
sync.WaitGroup
==============

What is it?
-----------
WaitGroup is a synchronization primitive that waits for a collection of
goroutines to finish executing. It's the standard way to wait for multiple
concurrent operations to complete.

Key Methods:
	- Add(delta int): Add delta to the WaitGroup counter
	- Done(): Decrement the counter by 1 (same as Add(-1))
	- Wait(): Block until counter becomes 0

How it works:
	1. Main goroutine calls Add(n) to set counter to n
	2. Spawn n goroutines, each calls Done() when finished
	3. Main goroutine calls Wait() to block until counter reaches 0

Benefits:
- Simple and safe goroutine synchronization
- Prevents main from exiting before goroutines finish
- No need for complex channel coordination
- Built-in and efficient

Use Cases:
- Waiting for multiple concurrent tasks to complete
- Parallel processing where you need all results
- Batch operations (file processing, API calls, etc.)
- Test scenarios with concurrent operations
- Worker pool coordination

Best Practices:
	1. Call Add() before starting goroutines (avoid race conditions)
	2. Call Done() with defer to ensure it's always called
	3. Don't pass WaitGroup by value (use pointer)
	4. Don't Wait() inside the same goroutine that calls Done()
	5. Counter must not go negative (will panic)
	6. Don't reuse WaitGroup until all Wait() calls return
	7. A WaitGroup is passed by reference because copying it creates independent counters, breaking synchronization.
	   (All goroutines must operate on the same shared state)
*/

// ============================================================================
// BASIC WAITGROUP USAGE
// ============================================================================

func wgWorker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // Ensure Done() is called even if panic occurs

	fmt.Printf("Worker %d: Starting\n", id)
	time.Sleep(time.Duration(id*100) * time.Millisecond)
	fmt.Printf("Worker %d: Finished\n", id)
}

func demoBasicWaitGroup() {
	fmt.Println("\n📌 Basic WaitGroup Usage")
	fmt.Println(strings.Repeat("-", 40))

	var wg sync.WaitGroup
	numWorkers := 5

	fmt.Printf("Starting %d workers...\n", numWorkers)

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1) // Increment counter before starting goroutine
		go wgWorker(i, &wg)
	}

	fmt.Println("Main: Waiting for workers to finish...")
	wg.Wait() // Block until all workers call Done()
	fmt.Println("✓ All workers finished!")
}

// ============================================================================
// WAITGROUP WITH ERRORS (Collecting Results)
// ============================================================================

type JobResult struct {
	ID     int
	Result string
	Error  error
}

func processJob(id int, wg *sync.WaitGroup, results chan<- JobResult) {
	defer wg.Done()

	fmt.Printf("Processing job %d...\n", id)
	time.Sleep(300 * time.Millisecond)

	// Simulate success/failure
	var result JobResult
	if id%3 == 0 {
		result = JobResult{ID: id, Result: "", Error: fmt.Errorf("job %d failed", id)}
	} else {
		result = JobResult{ID: id, Result: fmt.Sprintf("Job %d completed", id), Error: nil}
	}

	results <- result
}

func demoWaitGroupWithResults() {
	fmt.Println("\n📌 WaitGroup with Result Collection")
	fmt.Println(strings.Repeat("-", 40))

	var wg sync.WaitGroup
	numJobs := 6
	results := make(chan JobResult, numJobs)

	// Start jobs
	for i := 1; i <= numJobs; i++ {
		wg.Add(1)
		go processJob(i, &wg, results)
	}

	// Close results channel when all jobs are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	successCount := 0
	failCount := 0

	for result := range results {
		if result.Error != nil {
			fmt.Printf("  ✗ %v\n", result.Error)
			failCount++
		} else {
			fmt.Printf("  ✓ %s\n", result.Result)
			successCount++
		}
	}

	fmt.Printf("\nSummary: %d succeeded, %d failed\n", successCount, failCount)
}

// ============================================================================
// NESTED WAITGROUPS
// ============================================================================

func subTask(taskID, subID int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("    Subtask %d.%d executing\n", taskID, subID)
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("    Subtask %d.%d done\n", taskID, subID)
}

func task(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("  Task %d: Starting\n", id)

	// Nested WaitGroup for subtasks
	var subWg sync.WaitGroup
	numSubtasks := 3

	for i := 1; i <= numSubtasks; i++ {
		subWg.Add(1)
		go subTask(id, i, &subWg)
	}

	subWg.Wait() // Wait for all subtasks
	fmt.Printf("  Task %d: All subtasks completed\n", id)
}

func demoNestedWaitGroups() {
	fmt.Println("\n📌 Nested WaitGroups")
	fmt.Println(strings.Repeat("-", 40))

	var wg sync.WaitGroup
	numTasks := 2

	for i := 1; i <= numTasks; i++ {
		wg.Add(1)
		go task(i, &wg)
	}

	wg.Wait()
	fmt.Println("✓ All tasks and subtasks completed")
}

// ============================================================================
// WAITGROUP WITH DYNAMIC GOROUTINES
// ============================================================================

func dynamicWorker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Printf("Worker %d: No more jobs, exiting\n", id)
}

func demoDynamicGoroutines() {
	fmt.Println("\n📌 Dynamic Goroutines with WaitGroup")
	fmt.Println(strings.Repeat("-", 40))

	var wg sync.WaitGroup
	jobs := make(chan int, 10)
	numWorkers := 3

	// Start workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go dynamicWorker(i, jobs, &wg)
	}

	// Send jobs
	numJobs := 10
	fmt.Printf("Sending %d jobs to %d workers...\n", numJobs, numWorkers)
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // No more jobs

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("✓ All jobs processed")
}

// ============================================================================
// COMMON MISTAKE: Counter going negative (DEMO ONLY - DON'T DO THIS)
// ============================================================================

func demoWaitGroupCommonMistakes() {
	fmt.Println("\n📌 Common Mistakes to Avoid")
	fmt.Println(strings.Repeat("-", 40))

	fmt.Println("\n⚠ Mistake 1: Calling Add() inside goroutine")
	fmt.Println("   Problem: Race condition - Wait() might be called before Add()")
	fmt.Println("   Solution: Always call Add() before starting goroutine")

	fmt.Println("\n⚠ Mistake 2: Passing WaitGroup by value")
	fmt.Println("   Problem: Copy doesn't share state, Wait() never returns")
	fmt.Println("   Solution: Always pass *sync.WaitGroup (pointer)")

	fmt.Println("\n⚠ Mistake 3: Forgetting to call Done()")
	fmt.Println("   Problem: Wait() blocks forever (deadlock)")
	fmt.Println("   Solution: Use defer wg.Done() at start of function")

	fmt.Println("\n⚠ Mistake 4: Calling Done() more times than Add()")
	fmt.Println("   Problem: Counter goes negative, program panics")
	fmt.Println("   Solution: Ensure Add() count matches Done() calls")

	fmt.Println("\n✓ Following these practices prevents common bugs")
}

// ============================================================================
// PRACTICAL EXAMPLE: Parallel File Processing
// ============================================================================

type FileTask struct {
	Filename string
	Size     int
}

func processFile(file FileTask, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	fmt.Printf("  Processing %s (%d bytes)...\n", file.Filename, file.Size)
	time.Sleep(time.Duration(file.Size) * time.Millisecond)

	results <- fmt.Sprintf("%s processed successfully", file.Filename)
}

func demoWgPracticalExample() {
	fmt.Println("\n📌 Practical Example: Parallel File Processing")
	fmt.Println(strings.Repeat("-", 40))

	files := []FileTask{
		{"document1.txt", 100},
		{"image1.jpg", 200},
		{"video1.mp4", 300},
		{"document2.txt", 150},
		{"image2.jpg", 180},
	}

	var wg sync.WaitGroup
	results := make(chan string, len(files))

	fmt.Printf("Processing %d files in parallel...\n", len(files))
	start := time.Now()

	// Process all files concurrently
	for _, file := range files {
		wg.Add(1)
		go processFile(file, &wg, results)
	}

	// Wait and close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for result := range results {
		fmt.Printf("  ✓ %s\n", result)
	}

	duration := time.Since(start)
	fmt.Printf("\n✓ All files processed in %v\n", duration)
	fmt.Println("  (Would take ~930ms sequentially, but runs concurrently)")
}

// ============================================================================
// DEMONSTRATION
// ============================================================================

// DemoWaitGroup demonstrates all WaitGroup patterns
func DemoWaitGroup() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("sync.WaitGroup DEMONSTRATION")
	fmt.Println(strings.Repeat("=", 60) + "\n")

	demoBasicWaitGroup()
	demoWaitGroupWithResults()
	demoNestedWaitGroups()
	demoDynamicGoroutines()
	demoWaitGroupCommonMistakes()
	demoWgPracticalExample()

	fmt.Println("\n✅ WaitGroup demonstration completed!")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("- WaitGroup coordinates multiple goroutines")
	fmt.Println("- Add() before starting goroutine, Done() when finished")
	fmt.Println("- Wait() blocks until counter reaches zero")
	fmt.Println("- Always use defer wg.Done() for safety")
	fmt.Println("- Pass WaitGroup as pointer (*sync.WaitGroup)")
}
