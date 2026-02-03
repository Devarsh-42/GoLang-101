package main

import (
	"fmt"
	"sync"
	"time"
)

/*
sync.Mutex & sync.RWMutex
==========================

What is it?
-----------
Mutexes (Mutual Exclusion locks) protect shared data from concurrent access,
preventing race conditions. Go provides two types:

1. sync.Mutex: Basic mutual exclusion lock
   - Lock(): Acquire lock (blocks if already locked)
   - Unlock(): Release lock

2. sync.RWMutex: Read-Write lock
   - RLock(): Acquire read lock (multiple readers allowed)
   - RUnlock(): Release read lock
   - Lock(): Acquire write lock (exclusive access)
   - Unlock(): Release write lock

Benefits:
- Prevents race conditions
- Ensures data consistency
- Fine-grained control over critical sections
- RWMutex optimizes for read-heavy workloads

Use Cases:
- Protecting shared counters, maps, slices
- Cache implementations
- Configuration updates
- Database connection pools
- Any shared mutable state

Best Practices:
1. Keep critical sections small (minimize lock duration)
2. Use defer mu.Unlock() to ensure unlock happens
3. Don't copy mutexes (pass by pointer)
4. Avoid nested locks (can cause deadlock)
5. Use RWMutex when reads greatly outnumber writes
6. Consider channels for some scenarios (share by communicating)
7. Document what each mutex protects
*/

// ============================================================================
// BASIC MUTEX USAGE (Preventing Race Conditions)
// ============================================================================

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func demoBasicMutex() {
	fmt.Println("\n📌 Basic Mutex Usage")
	fmt.Println("-" * 40)

	counter := &Counter{}
	var wg sync.WaitGroup

	// Without mutex, this would cause race condition
	numGoroutines := 100
	incrementsPerGoroutine := 100

	fmt.Printf("Starting %d goroutines, each incrementing %d times\n",
		numGoroutines, incrementsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				counter.Increment()
			}
		}(i)
	}

	wg.Wait()

	expected := numGoroutines * incrementsPerGoroutine
	actual := counter.Value()

	fmt.Printf("\nExpected: %d\n", expected)
	fmt.Printf("Actual:   %d\n", actual)

	if expected == actual {
		fmt.Println("✓ Mutex prevented race condition!")
	} else {
		fmt.Println("✗ Race condition occurred!")
	}
}

// ============================================================================
// RACE CONDITION DEMONSTRATION (Without Mutex)
// ============================================================================

func demoRaceCondition() {
	fmt.Println("\n📌 Race Condition (Without Mutex)")
	fmt.Println("-" * 40)

	var counter int // No mutex protection
	var wg sync.WaitGroup

	numGoroutines := 100

	fmt.Println("⚠  Running concurrent increments WITHOUT mutex...")
	fmt.Println("   (Run with 'go run -race' to detect race condition)")

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter++ // RACE CONDITION!
			}
		}()
	}

	wg.Wait()

	expected := numGoroutines * 100
	fmt.Printf("\nExpected: %d\n", expected)
	fmt.Printf("Actual:   %d\n", counter)
	fmt.Println("✗ Results may vary due to race condition!")
}

// ============================================================================
// RWMutex (Read-Write Lock)
// ============================================================================

type Cache struct {
	mu     sync.RWMutex
	data   map[string]string
	reads  int
	writes int
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock() // Multiple readers can acquire read lock
	defer c.mu.RUnlock()

	c.reads++
	value, ok := c.data[key]
	return value, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock() // Write lock is exclusive
	defer c.mu.Unlock()

	c.writes++
	c.data[key] = value
}

func (c *Cache) Stats() (reads, writes int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.reads, c.writes
}

func demoRWMutex() {
	fmt.Println("\n📌 RWMutex (Read-Write Lock)")
	fmt.Println("-" * 40)

	cache := NewCache()
	var wg sync.WaitGroup

	// Pre-populate cache
	for i := 1; i <= 10; i++ {
		cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	// Many readers (90% of operations)
	numReaders := 90
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", (id%10)+1)
			value, ok := cache.Get(key)
			if ok {
				_ = value // Use value
			}
		}(i)
	}

	// Few writers (10% of operations)
	numWriters := 10
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cache.Set(fmt.Sprintf("newkey%d", id), fmt.Sprintf("newvalue%d", id))
			time.Sleep(10 * time.Millisecond) // Simulate slower write
		}(i)
	}

	wg.Wait()

	reads, writes := cache.Stats()
	fmt.Printf("\nCache Statistics:\n")
	fmt.Printf("  Reads:  %d (multiple concurrent readers allowed)\n", reads)
	fmt.Printf("  Writes: %d (exclusive access required)\n", writes)
	fmt.Println("\n✓ RWMutex allows multiple concurrent readers!")
	fmt.Println("  (Better performance for read-heavy workloads)")
}

// ============================================================================
// COMMON PATTERN: Protecting a Map
// ============================================================================

type SafeMap struct {
	mu   sync.RWMutex
	data map[string]int
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]int),
	}
}

func (sm *SafeMap) Set(key string, value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

func (sm *SafeMap) Get(key string) (int, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	val, ok := sm.data[key]
	return val, ok
}

func (sm *SafeMap) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.data, key)
}

func (sm *SafeMap) Len() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.data)
}

func demoSafeMap() {
	fmt.Println("\n📌 Thread-Safe Map")
	fmt.Println("-" * 40)

	safeMap := NewSafeMap()
	var wg sync.WaitGroup

	// Concurrent writes
	fmt.Println("Performing concurrent map operations...")
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			safeMap.Set(fmt.Sprintf("key%d", id), id*10)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", id%25)
			if val, ok := safeMap.Get(key); ok {
				_ = val
			}
		}(i)
	}

	wg.Wait()

	fmt.Printf("\nFinal map size: %d entries\n", safeMap.Len())
	fmt.Println("✓ All concurrent operations completed safely!")
}

// ============================================================================
// PERFORMANCE COMPARISON: Mutex vs RWMutex
// ============================================================================

func demoPerformanceComparison() {
	fmt.Println("\n📌 Performance: Mutex vs RWMutex")
	fmt.Println("-" * 40)

	// Test data
	data := make(map[string]string)
	for i := 0; i < 100; i++ {
		data[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
	}

	// Test with regular Mutex
	fmt.Println("\nTest 1: Using sync.Mutex (all operations exclusive)")
	var mu1 sync.Mutex
	var wg1 sync.WaitGroup

	start1 := time.Now()
	for i := 0; i < 1000; i++ {
		wg1.Add(1)
		go func(id int) {
			defer wg1.Done()
			mu1.Lock()
			_ = data[fmt.Sprintf("key%d", id%100)] // Read operation
			mu1.Unlock()
		}(i)
	}
	wg1.Wait()
	duration1 := time.Since(start1)

	// Test with RWMutex
	fmt.Println("Test 2: Using sync.RWMutex (reads concurrent)")
	var mu2 sync.RWMutex
	var wg2 sync.WaitGroup

	start2 := time.Now()
	for i := 0; i < 1000; i++ {
		wg2.Add(1)
		go func(id int) {
			defer wg2.Done()
			mu2.RLock() // Read lock
			_ = data[fmt.Sprintf("key%d", id%100)]
			mu2.RUnlock()
		}(i)
	}
	wg2.Wait()
	duration2 := time.Since(start2)

	fmt.Printf("\nResults:\n")
	fmt.Printf("  Mutex:   %v\n", duration1)
	fmt.Printf("  RWMutex: %v\n", duration2)

	if duration2 < duration1 {
		speedup := float64(duration1) / float64(duration2)
		fmt.Printf("  ✓ RWMutex is %.2fx faster for read-heavy workload!\n", speedup)
	}
}

// ============================================================================
// COMMON MISTAKES
// ============================================================================

func demoCommonMistakes() {
	fmt.Println("\n📌 Common Mistakes to Avoid")
	fmt.Println("-" * 40)

	fmt.Println("\n⚠ Mistake 1: Forgetting to unlock")
	fmt.Println("   Problem: Deadlock - other goroutines wait forever")
	fmt.Println("   Solution: Always use 'defer mu.Unlock()'")

	fmt.Println("\n⚠ Mistake 2: Copying mutex by value")
	fmt.Println("   Problem: Copy doesn't share state, no synchronization")
	fmt.Println("   Solution: Always pass mutex by pointer or embed in struct")

	fmt.Println("\n⚠ Mistake 3: Locking inside locked section (deadlock)")
	fmt.Println("   Problem: Same goroutine tries to lock already-locked mutex")
	fmt.Println("   Solution: Avoid nested locks, restructure code")

	fmt.Println("\n⚠ Mistake 4: Holding lock too long")
	fmt.Println("   Problem: Poor performance, increased contention")
	fmt.Println("   Solution: Keep critical sections small and fast")

	fmt.Println("\n⚠ Mistake 5: Not using RWMutex for read-heavy workloads")
	fmt.Println("   Problem: Unnecessary serialization of reads")
	fmt.Println("   Solution: Use RWMutex when reads >> writes")

	fmt.Println("\n✓ Follow these practices for safe concurrent code")
}

// ============================================================================
// DEMONSTRATION
// ============================================================================

// DemoMutex demonstrates all mutex patterns
func DemoMutex() {
	fmt.Println("\n" + "="*60)
	fmt.Println("sync.Mutex & sync.RWMutex DEMONSTRATION")
	fmt.Println("="*60 + "\n")

	demoBasicMutex()
	demoRaceCondition()
	demoRWMutex()
	demoSafeMap()
	demoPerformanceComparison()
	demoCommonMistakes()

	fmt.Println("\n✅ Mutex demonstration completed!")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("- Mutex provides mutual exclusion for shared data")
	fmt.Println("- Always use defer mu.Unlock() for safety")
	fmt.Println("- RWMutex allows multiple concurrent readers")
	fmt.Println("- Keep critical sections small")
	fmt.Println("- Use 'go run -race' to detect race conditions")
}
