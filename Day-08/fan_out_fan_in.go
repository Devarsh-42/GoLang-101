package main

import (
	"fmt"
	"sync"
	"time"
)

/*
FAN-OUT / FAN-IN PATTERN
========================

What is it?
-----------
Fan-out/Fan-in is a concurrency pattern where:
- FAN-OUT: Multiple goroutines are spawned to handle parts of a task concurrently
- FAN-IN: Multiple channels are merged into a single channel to collect results

This pattern is useful for parallelizing work and then aggregating results.

Key Concepts:
- Fan-out: Distribute work across multiple goroutines
- Fan-in: Collect results from multiple goroutines into one channel
- Pipeline: Chain of stages connected by channels

Benefits:
- Maximizes parallelism for independent tasks
- Efficient resource utilization
- Scalable processing

Use Cases:
- Parallel data processing (map-reduce style operations)
- Concurrent API calls to multiple services
- Distributed search across multiple data sources
- Image processing (multiple filters applied in parallel)
- Log aggregation from multiple sources

Best Practices:
1. Ensure work is truly independent (no shared state)
2. Use buffered channels to prevent blocking
3. Close channels properly to avoid goroutine leaks
4. Use sync.WaitGroup to coordinate goroutine completion
5. Consider context for cancellation
*/

// ============================================================================
// FAN-OUT: Distribute work to multiple workers
// ============================================================================

// squareWorker squares numbers from input channel
func squareWorker(id int, numbers <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for num := range numbers {
		fmt.Printf("Worker %d squaring %d\n", id, num)
		time.Sleep(100 * time.Millisecond) // Simulate work
		results <- num * num
	}
}

// fanOutSquare demonstrates fan-out pattern
func fanOutSquare(numbers []int, numWorkers int) <-chan int {
	// Create channels
	input := make(chan int, len(numbers))
	output := make(chan int, len(numbers))

	var wg sync.WaitGroup

	// FAN-OUT: Start multiple workers
	fmt.Printf("\n🔀 Fan-out: Starting %d workers\n", numWorkers)
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go squareWorker(i, input, output, &wg)
	}

	// Send work to workers
	go func() {
		for _, num := range numbers {
			input <- num
		}
		close(input)
	}()

	// Close output channel when all workers are done
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// ============================================================================
// FAN-IN: Merge multiple channels into one
// ============================================================================

// merge combines multiple channels into a single channel (fan-in)
func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	merged := make(chan int)

	// Start a goroutine for each input channel
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			merged <- n
		}
	}

	// FAN-IN: Collect from all channels
	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}

	// Close merged channel when all inputs are exhausted
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

// generator creates a channel and sends numbers
func generator(start, count int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := start; i < start+count; i++ {
			out <- i
			time.Sleep(50 * time.Millisecond)
		}
	}()
	return out
}

// processor processes numbers (squares them)
func processor(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			time.Sleep(100 * time.Millisecond) // Simulate processing
			out <- n * n
		}
	}()
	return out
}

// ============================================================================
// DEMONSTRATION
// ============================================================================

// DemoFanOutFanIn demonstrates both fan-out and fan-in patterns
func DemoFanOutFanIn() {
	fmt.Println("\n" + "="*60)
	fmt.Println("FAN-OUT / FAN-IN DEMONSTRATION")
	fmt.Println("="*60 + "\n")

	// ========================================
	// Part 1: Simple Fan-Out Example
	// ========================================
	fmt.Println("📤 Part 1: Fan-Out Pattern")
	fmt.Println("-" * 40)

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("Input: %v\n", numbers)

	squared := fanOutSquare(numbers, 3) // 3 workers

	fmt.Println("\n📥 Collecting results:")
	var results []int
	for result := range squared {
		results = append(results, result)
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Printf("Results: %v\n", results)

	// ========================================
	// Part 2: Fan-In Pattern
	// ========================================
	fmt.Println("\n📤 Part 2: Fan-In Pattern")
	fmt.Println("-" * 40)

	// Create multiple generators (fan-out)
	gen1 := generator(1, 3)  // generates 1, 2, 3
	gen2 := generator(10, 3) // generates 10, 11, 12
	gen3 := generator(20, 3) // generates 20, 21, 22

	fmt.Println("Creating 3 generators producing numbers concurrently...")

	// Merge all generators (fan-in)
	merged := merge(gen1, gen2, gen3)

	fmt.Println("\n📥 Receiving merged output:")
	var mergedResults []int
	for num := range merged {
		mergedResults = append(mergedResults, num)
		fmt.Printf("Received: %d\n", num)
	}
	fmt.Printf("All merged results: %v\n", mergedResults)

	// ========================================
	// Part 3: Complete Pipeline (Fan-out + Fan-in)
	// ========================================
	fmt.Println("\n📤 Part 3: Complete Pipeline")
	fmt.Println("-" * 40)

	// Stage 1: Generate numbers
	input := generator(1, 5)

	// Stage 2: Fan-out - multiple processors
	fmt.Println("Fan-out: Starting 3 processors in parallel...")
	proc1 := processor(input)
	proc2 := processor(input)
	proc3 := processor(input)

	// Stage 3: Fan-in - merge results
	fmt.Println("Fan-in: Merging results from all processors...")
	finalOutput := merge(proc1, proc2, proc3)

	// Collect final results
	fmt.Println("\n📥 Final pipeline output:")
	var pipelineResults []int
	for result := range finalOutput {
		pipelineResults = append(pipelineResults, result)
		fmt.Printf("Final result: %d\n", result)
	}

	fmt.Println("\n✅ Fan-out/Fan-in demonstration completed!")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("- Fan-out distributes work to multiple goroutines")
	fmt.Println("- Fan-in merges results from multiple channels")
	fmt.Println("- Can be combined to create efficient pipelines")
	fmt.Println("- Great for parallel processing of independent tasks")
}
