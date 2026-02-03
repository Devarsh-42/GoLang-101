package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/*
ATOMIC OPERATIONS
=================

What is it?
-----------
Atomic operations are CPU-level operations that complete in a single,
indivisible step without the possibility of interference from other goroutines.
The sync/atomic package provides low-level atomic memory primitives.

Common Atomic Operations:
- atomic.AddInt32/Int64: Atomic addition
- atomic.LoadInt32/Int64: Atomic read
- atomic.StoreInt32/Int64: Atomic write
- atomic.SwapInt32/Int64: Atomic exchange
- atomic.CompareAndSwapInt32/Int64: Compare-and-swap (CAS)
- atomic.Value: Atomic operations for any value type

Benefits:
- Lock-free synchronization (no mutex overhead)
- Better performance than mutexes for simple operations
- No risk of deadlock
- Hardware-level guarantees

When to Use:
- Simple counters, flags, or single values
- Performance-critical code paths
- Lock-free data structures
- Metrics and statistics collection

When NOT to Use:
- Complex operations on multiple values
- Operations requiring consistency across multiple fields
- When code clarity is more important than performance

Best Practices:
1. Use atomic operations for simple scenarios only
2. Prefer mutexes for complex state
3. Document why atomic is used instead of mutex
4. Use atomic.Value for non-numeric types
5. Understand memory ordering guarantees
*/

// ============================================================================
// BASIC ATOMIC OPERATIONS
// ============================================================================

func demoBasicAtomic() {
	fmt.Println("\n📌 Basic Atomic Operations")
	fmt.Println("-" * 40)

	var counter int64 = 0
	var wg sync.WaitGroup

	numGoroutines := 100
	incrementsPerGoroutine := 100

	fmt.Printf("Starting %d goroutines, each incrementing %d times\n",
		numGoroutines, incrementsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				atomic.AddInt64(&counter, 1) // Atomic increment
			}
		}()
	}

	wg.Wait()

	expected := int64(numGoroutines * incrementsPerGoroutine)
	actual := atomic.LoadInt64(&counter) // Atomic read

	fmt.Printf("\nExpected: %d\n", expected)
	fmt.Printf("Actual:   %d\n", actual)

	if expected == actual {
		fmt.Println("✓ Atomic operations prevented race condition!")
	}
}

// ============================================================================
// ATOMIC VS MUTEX PERFORMANCE
// ============================================================================

func demoAtomicVsMutex() {
	fmt.Println("\n📌 Performance: Atomic vs Mutex")
	fmt.Println("-" * 40)

	iterations := 100000

	// Test 1: Atomic operations
	var atomicCounter int64
	start1 := time.Now()

	var wg1 sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			for j := 0; j < iterations/10; j++ {
				atomic.AddInt64(&atomicCounter, 1)
			}
		}()
	}
	wg1.Wait()
	duration1 := time.Since(start1)

	// Test 2: Mutex-based counter
	type MutexCounter struct {
		mu    sync.Mutex
		value int64
	}
	mutexCounter := &MutexCounter{}

	start2 := time.Now()
	var wg2 sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			for j := 0; j < iterations/10; j++ {
				mutexCounter.mu.Lock()
				mutexCounter.value++
				mutexCounter.mu.Unlock()
			}
		}()
	}
	wg2.Wait()
	duration2 := time.Since(start2)

	fmt.Printf("\nResults (%d operations):\n", iterations)
	fmt.Printf("  Atomic: %v (value: %d)\n", duration1, atomic.LoadInt64(&atomicCounter))
	fmt.Printf("  Mutex:  %v (value: %d)\n", duration2, mutexCounter.value)

	if duration1 < duration2 {
		speedup := float64(duration2) / float64(duration1)
		fmt.Printf("  ✓ Atomic is %.2fx faster!\n", speedup)
	}
}

// ============================================================================
// ATOMIC LOAD AND STORE
// ============================================================================

func demoLoadStore() {
	fmt.Println("\n📌 Atomic Load and Store")
	fmt.Println("-" * 40)

	var config int64 = 0
	var wg sync.WaitGroup

	// Writer: Updates config every 100ms
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 5; i++ {
			time.Sleep(100 * time.Millisecond)
			atomic.StoreInt64(&config, int64(i))
			fmt.Printf("Writer: Updated config to %d\n", i)
		}
	}()

	// Readers: Read config concurrently
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				value := atomic.LoadInt64(&config)
				fmt.Printf("  Reader %d: config = %d\n", id, value)
				time.Sleep(50 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("✓ Atomic load/store ensures consistent reads")
}

// ============================================================================
// COMPARE-AND-SWAP (CAS)
// ============================================================================

func demoCompareAndSwap() {
	fmt.Println("\n📌 Compare-and-Swap (CAS)")
	fmt.Println("-" * 40)

	var value int64 = 100
	var wg sync.WaitGroup

	fmt.Println("Multiple goroutines trying to update value from 100 to their ID")
	fmt.Printf("Initial value: %d\n\n", value)

	// Multiple goroutines try to CAS
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Try to swap value from 100 to id
			swapped := atomic.CompareAndSwapInt64(&value, 100, int64(id))

			if swapped {
				fmt.Printf("  ✓ Goroutine %d: Successfully swapped value to %d\n", id, id)
			} else {
				currentValue := atomic.LoadInt64(&value)
				fmt.Printf("  ✗ Goroutine %d: Failed (value is now %d, not 100)\n", id, currentValue)
			}
		}(i)
	}

	wg.Wait()
	finalValue := atomic.LoadInt64(&value)
	fmt.Printf("\nFinal value: %d\n", finalValue)
	fmt.Println("✓ Only one goroutine succeeded (atomic guarantee)")
}

// ============================================================================
// ATOMIC.VALUE (For Any Type)
// ============================================================================

type Config struct {
	MaxConnections int
	Timeout        time.Duration
	EnableLogging  bool
}

func demoAtomicValue() {
	fmt.Println("\n📌 atomic.Value (Any Type)")
	fmt.Println("-" * 40)

	var atomicConfig atomic.Value

	// Initial config
	initialConfig := Config{
		MaxConnections: 10,
		Timeout:        5 * time.Second,
		EnableLogging:  false,
	}
	atomicConfig.Store(initialConfig)

	var wg sync.WaitGroup

	// Readers
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				cfg := atomicConfig.Load().(Config)
				fmt.Printf("  Reader %d: MaxConn=%d, Timeout=%v, Logging=%v\n",
					id, cfg.MaxConnections, cfg.Timeout, cfg.EnableLogging)
				time.Sleep(100 * time.Millisecond)
			}
		}(i)
	}

	// Writer: Update config after 250ms
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(250 * time.Millisecond)

		newConfig := Config{
			MaxConnections: 100,
			Timeout:        10 * time.Second,
			EnableLogging:  true,
		}
		atomicConfig.Store(newConfig)
		fmt.Println("\n  Writer: Config updated!\n")
	}()

	wg.Wait()
	fmt.Println("✓ atomic.Value safely stores/loads any type")
}

// ============================================================================
// PRACTICAL EXAMPLE: Thread-Safe Statistics
// ============================================================================

type Statistics struct {
	requests  int64
	errors    int64
	totalTime int64
}

func (s *Statistics) RecordRequest(duration time.Duration, hasError bool) {
	atomic.AddInt64(&s.requests, 1)
	atomic.AddInt64(&s.totalTime, int64(duration))

	if hasError {
		atomic.AddInt64(&s.errors, 1)
	}
}

func (s *Statistics) GetStats() (requests, errors int64, avgTime time.Duration) {
	req := atomic.LoadInt64(&s.requests)
	err := atomic.LoadInt64(&s.errors)
	total := atomic.LoadInt64(&s.totalTime)

	if req > 0 {
		avgTime = time.Duration(total / req)
	}

	return req, err, avgTime
}

func demoPracticalExample() {
	fmt.Println("\n📌 Practical Example: Thread-Safe Statistics")
	fmt.Println("-" * 40)

	stats := &Statistics{}
	var wg sync.WaitGroup

	// Simulate multiple goroutines recording statistics
	fmt.Println("Simulating concurrent request processing...")

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Simulate request
			start := time.Now()
			time.Sleep(time.Duration(id%10) * time.Millisecond)
			duration := time.Since(start)

			// Record with 20% error rate
			hasError := id%5 == 0
			stats.RecordRequest(duration, hasError)
		}(i)
	}

	wg.Wait()

	// Get final statistics
	requests, errors, avgTime := stats.GetStats()

	fmt.Printf("\nStatistics:\n")
	fmt.Printf("  Total Requests: %d\n", requests)
	fmt.Printf("  Errors:         %d (%.1f%%)\n", errors, float64(errors)/float64(requests)*100)
	fmt.Printf("  Avg Time:       %v\n", avgTime)
	fmt.Println("\n✓ Lock-free statistics collection!")
}

// ============================================================================
// WHEN TO USE ATOMIC VS MUTEX
// ============================================================================

func demoWhenToUse() {
	fmt.Println("\n📌 When to Use Atomic vs Mutex")
	fmt.Println("-" * 40)

	fmt.Println("\n✅ Use ATOMIC when:")
	fmt.Println("  • Single value operations (counter, flag)")
	fmt.Println("  • Performance is critical")
	fmt.Println("  • Simple increment/decrement/read/write")
	fmt.Println("  • Lock-free algorithms")

	fmt.Println("\n✅ Use MUTEX when:")
	fmt.Println("  • Multiple related values need consistency")
	fmt.Println("  • Complex operations on shared state")
	fmt.Println("  • Code clarity is more important than performance")
	fmt.Println("  • Operations span multiple lines")

	fmt.Println("\n📊 Example comparison:")
	fmt.Println("\n  // Atomic (simple counter)")
	fmt.Println("  atomic.AddInt64(&counter, 1)")

	fmt.Println("\n  // Mutex (complex state)")
	fmt.Println("  mu.Lock()")
	fmt.Println("  account.balance -= amount")
	fmt.Println("  account.transactions++")
	fmt.Println("  account.lastUpdate = time.Now()")
	fmt.Println("  mu.Unlock()")
}

// ============================================================================
// DEMONSTRATION
// ============================================================================

// DemoAtomic demonstrates all atomic patterns
func DemoAtomic() {
	fmt.Println("\n" + "="*60)
	fmt.Println("ATOMIC OPERATIONS DEMONSTRATION")
	fmt.Println("="*60 + "\n")

	demoBasicAtomic()
	demoAtomicVsMutex()
	demoLoadStore()
	demoCompareAndSwap()
	demoAtomicValue()
	demoPracticalExample()
	demoWhenToUse()

	fmt.Println("\n✅ Atomic operations demonstration completed!")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("- Atomic operations are lock-free and fast")
	fmt.Println("- Use for simple single-value operations")
	fmt.Println("- Better performance than mutex for counters")
	fmt.Println("- atomic.Value works with any type")
	fmt.Println("- CAS enables lock-free algorithms")
	fmt.Println("- Prefer mutex for complex state management")
}
