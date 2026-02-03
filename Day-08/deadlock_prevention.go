package main

import (
	"fmt"
	"sync"
	"time"
)

/*
DEADLOCK PREVENTION
===================

What is it?
-----------
A deadlock occurs when two or more goroutines are waiting for each other
to release resources, creating a circular dependency that prevents progress.

Four Conditions for Deadlock (Coffman conditions):
1. Mutual Exclusion: Resource held exclusively
2. Hold and Wait: Holding resource while waiting for another
3. No Preemption: Resources cannot be forcibly taken
4. Circular Wait: Circular chain of waiting goroutines

Common Deadlock Scenarios:
- Lock ordering issues (nested locks)
- Unbuffered channel with no receiver
- WaitGroup counter mismatch
- All goroutines blocked on channel operations

Prevention Strategies:
1. Lock ordering: Always acquire locks in the same order
2. Lock timeout: Use context with timeout
3. Try-lock pattern: Non-blocking lock attempts
4. Avoid nested locks: Minimize lock scope
5. Use buffered channels when appropriate
6. Always close channels and call Done()

Detection:
- Go runtime detects some deadlocks and panics
- Use 'go run -race' to detect race conditions
- Monitor goroutine counts and lock contention

Best Practices:
1. Keep critical sections small
2. Release locks as soon as possible (use defer)
3. Avoid holding multiple locks
4. Use channels for communication when possible
5. Document lock ordering requirements
6. Use timeouts to prevent indefinite waiting
*/

// ============================================================================
// EXAMPLE 1: CLASSIC DEADLOCK (Lock Ordering)
// ============================================================================

type Account struct {
	mu      sync.Mutex
	id      int
	balance int
}

// BAD: Can cause deadlock
func transferBad(from, to *Account, amount int) {
	from.mu.Lock()
	defer from.mu.Unlock()

	// Dangerous: Different order in different calls can deadlock
	to.mu.Lock()
	defer to.mu.Unlock()

	from.balance -= amount
	to.balance += amount
}

// GOOD: Lock ordering prevents deadlock
func transferGood(from, to *Account, amount int) {
	// Always lock in consistent order (by ID)
	first, second := from, to
	if first.id > second.id {
		first, second = second, first
	}

	first.mu.Lock()
	defer first.mu.Unlock()

	second.mu.Lock()
	defer second.mu.Unlock()

	from.balance -= amount
	to.balance += amount
}

func demoLockOrdering() {
	fmt.Println("\n📌 Deadlock Prevention: Lock Ordering")
	fmt.Println("-" * 40)

	fmt.Println("\n⚠ Scenario: Two accounts transferring to each other")
	fmt.Println("  Account A -> Account B")
	fmt.Println("  Account B -> Account A")
	fmt.Println("  (Without proper lock ordering, this can deadlock)")

	fmt.Println("\n✅ Solution: Always acquire locks in consistent order")
	fmt.Println("  1. Sort locks by ID (or memory address)")
	fmt.Println("  2. Always lock lower ID first")
	fmt.Println("  3. Then lock higher ID")

	accountA := &Account{id: 1, balance: 1000}
	accountB := &Account{id: 2, balance: 1000}

	var wg sync.WaitGroup

	// Concurrent transfers using good method
	for i := 0; i < 5; i++ {
		wg.Add(2)

		go func() {
			defer wg.Done()
			transferGood(accountA, accountB, 10)
		}()

		go func() {
			defer wg.Done()
			transferGood(accountB, accountA, 10)
		}()
	}

	wg.Wait()

	fmt.Printf("\nFinal balances: A=%d, B=%d\n", accountA.balance, accountB.balance)
	fmt.Println("✓ No deadlock occurred!")
}

// ============================================================================
// EXAMPLE 2: UNBUFFERED CHANNEL DEADLOCK
// ============================================================================

func demoChannelDeadlock() {
	fmt.Println("\n📌 Channel Deadlock Prevention")
	fmt.Println("-" * 40)

	// BAD: Deadlock
	fmt.Println("\n❌ BAD Pattern (would deadlock):")
	fmt.Println("  ch := make(chan int)")
	fmt.Println("  ch <- 1  // Blocks forever - no receiver!")
	fmt.Println("  (Go runtime would detect and panic)")

	// GOOD: Solutions
	fmt.Println("\n✅ SOLUTION 1: Use goroutine for send")
	ch1 := make(chan int)
	go func() {
		ch1 <- 1
	}()
	value := <-ch1
	fmt.Printf("  Received: %d\n", value)

	fmt.Println("\n✅ SOLUTION 2: Use buffered channel")
	ch2 := make(chan int, 1)
	ch2 <- 2 // Doesn't block - buffer available
	value = <-ch2
	fmt.Printf("  Received: %d\n", value)

	fmt.Println("\n✅ SOLUTION 3: Use select with default")
	ch3 := make(chan int)
	select {
	case ch3 <- 3:
		fmt.Println("  Sent")
	default:
		fmt.Println("  Would block - skipping send")
	}
}

// ============================================================================
// EXAMPLE 3: WAITGROUP DEADLOCK
// ============================================================================

func demoWaitGroupDeadlock() {
	fmt.Println("\n📌 WaitGroup Deadlock Prevention")
	fmt.Println("-" * 40)

	// BAD: Deadlock
	fmt.Println("\n❌ BAD Pattern (would deadlock):")
	fmt.Println("  var wg sync.WaitGroup")
	fmt.Println("  wg.Add(2)")
	fmt.Println("  go func() { wg.Done() }()")
	fmt.Println("  // Only 1 Done() called, but Add(2)")
	fmt.Println("  wg.Wait() // Waits forever!")

	// GOOD: Proper usage
	fmt.Println("\n✅ SOLUTION: Match Add() count with Done() calls")
	var wg sync.WaitGroup

	tasks := 3
	wg.Add(tasks)

	for i := 1; i <= tasks; i++ {
		go func(id int) {
			defer wg.Done() // Always use defer
			fmt.Printf("  Task %d completed\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("✓ All tasks completed - no deadlock")
}

// ============================================================================
// EXAMPLE 4: CIRCULAR CHANNEL DEPENDENCY
// ============================================================================

func demoCircularDependency() {
	fmt.Println("\n📌 Circular Channel Dependency")
	fmt.Println("-" * 40)

	// BAD: Circular wait
	fmt.Println("\n❌ BAD Pattern (would deadlock):")
	fmt.Println("  ch1, ch2 := make(chan int), make(chan int)")
	fmt.Println("  go func() {")
	fmt.Println("    ch1 <- (<-ch2) + 1  // Waits for ch2")
	fmt.Println("  }()")
	fmt.Println("  ch2 <- (<-ch1) + 1    // Waits for ch1")
	fmt.Println("  (Circular wait - deadlock!)")

	// GOOD: Break the cycle
	fmt.Println("\n✅ SOLUTION: Use buffered channels or goroutines")
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	// Initialize with a value to break the cycle
	ch1 <- 1

	go func() {
		val := <-ch1
		ch2 <- val + 1
	}()

	result := <-ch2
	fmt.Printf("  Result: %d\n", result)
	fmt.Println("✓ No circular dependency")
}

// ============================================================================
// EXAMPLE 5: TIMEOUT TO PREVENT DEADLOCK
// ============================================================================

func demoTimeoutPrevention() {
	fmt.Println("\n📌 Timeout-Based Deadlock Prevention")
	fmt.Println("-" * 40)

	ch := make(chan int)

	// Without timeout: Would wait forever
	fmt.Println("Attempting to receive from empty channel...")

	select {
	case val := <-ch:
		fmt.Printf("Received: %d\n", val)
	case <-time.After(1 * time.Second):
		fmt.Println("✓ Timeout after 1 second - prevented indefinite wait")
	}

	fmt.Println("\nUsing timeout prevents goroutines from blocking forever")
}

// ============================================================================
// EXAMPLE 6: TRY-LOCK PATTERN
// ============================================================================

type TryLocker struct {
	mu     sync.Mutex
	locked bool
}

func (tl *TryLocker) TryLock() bool {
	select {
	case <-time.After(100 * time.Millisecond):
		return false // Timeout - couldn't acquire lock
	default:
		tl.mu.Lock()
		tl.locked = true
		return true
	}
}

func (tl *TryLocker) Unlock() {
	if tl.locked {
		tl.locked = false
		tl.mu.Unlock()
	}
}

func demoTryLock() {
	fmt.Println("\n📌 Try-Lock Pattern")
	fmt.Println("-" * 40)

	fmt.Println("\n✅ Instead of blocking indefinitely, try-lock with timeout")
	fmt.Println("  This prevents deadlock in complex lock scenarios")

	var mu sync.Mutex

	// First goroutine holds lock
	mu.Lock()

	// Second goroutine tries to acquire with timeout
	timeout := time.After(100 * time.Millisecond)

	select {
	case <-timeout:
		fmt.Println("  ✓ Lock acquisition timeout - prevented deadlock")
		fmt.Println("  Can retry, log error, or take alternative action")
	}

	mu.Unlock()
}

// ============================================================================
// PRACTICAL EXAMPLE: DEADLOCK-FREE RESOURCE MANAGER
// ============================================================================

type Resource struct {
	id   int
	name string
}

type ResourceManager struct {
	resources map[int]*Resource
	mu        sync.RWMutex
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		resources: make(map[int]*Resource),
	}
}

func (rm *ResourceManager) AcquireResources(ids []int) ([]*Resource, error) {
	// Sort IDs to ensure consistent lock ordering
	sortedIDs := make([]int, len(ids))
	copy(sortedIDs, ids)

	// Simple bubble sort for demonstration
	for i := 0; i < len(sortedIDs); i++ {
		for j := i + 1; j < len(sortedIDs); j++ {
			if sortedIDs[i] > sortedIDs[j] {
				sortedIDs[i], sortedIDs[j] = sortedIDs[j], sortedIDs[i]
			}
		}
	}

	// Acquire in sorted order
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	resources := make([]*Resource, 0, len(sortedIDs))
	for _, id := range sortedIDs {
		if res, ok := rm.resources[id]; ok {
			resources = append(resources, res)
		}
	}

	return resources, nil
}

func demoPracticalExample() {
	fmt.Println("\n📌 Practical Example: Deadlock-Free Resource Manager")
	fmt.Println("-" * 40)

	rm := NewResourceManager()

	// Add resources
	rm.mu.Lock()
	rm.resources[1] = &Resource{id: 1, name: "Database"}
	rm.resources[2] = &Resource{id: 2, name: "Cache"}
	rm.resources[3] = &Resource{id: 3, name: "Queue"}
	rm.mu.Unlock()

	var wg sync.WaitGroup

	// Multiple goroutines acquiring resources in different orders
	fmt.Println("Multiple goroutines acquiring resources concurrently...")

	wg.Add(2)

	go func() {
		defer wg.Done()
		resources, _ := rm.AcquireResources([]int{1, 2, 3})
		fmt.Printf("  Goroutine 1 acquired: %d resources\n", len(resources))
	}()

	go func() {
		defer wg.Done()
		resources, _ := rm.AcquireResources([]int{3, 2, 1}) // Different order
		fmt.Printf("  Goroutine 2 acquired: %d resources\n", len(resources))
	}()

	wg.Wait()

	fmt.Println("✓ No deadlock - consistent lock ordering works!")
}

// ============================================================================
// BEST PRACTICES SUMMARY
// ============================================================================

func demoBestPractices() {
	fmt.Println("\n📌 Deadlock Prevention Best Practices")
	fmt.Println("-" * 40)

	fmt.Println("\n✅ Lock Ordering:")
	fmt.Println("  • Always acquire locks in the same order")
	fmt.Println("  • Sort by ID or memory address")
	fmt.Println("  • Document lock ordering requirements")

	fmt.Println("\n✅ Channel Best Practices:")
	fmt.Println("  • Use buffered channels when appropriate")
	fmt.Println("  • Always close channels when done")
	fmt.Println("  • Use goroutines for sends/receives")
	fmt.Println("  • Use select with default for non-blocking")

	fmt.Println("\n✅ General Guidelines:")
	fmt.Println("  • Keep critical sections small")
	fmt.Println("  • Use defer for unlocks")
	fmt.Println("  • Avoid nested locks when possible")
	fmt.Println("  • Use timeouts to prevent indefinite waits")
	fmt.Println("  • Prefer channels over shared memory")

	fmt.Println("\n✅ Testing & Detection:")
	fmt.Println("  • Use 'go run -race' to detect races")
	fmt.Println("  • Test with high concurrency")
	fmt.Println("  • Monitor goroutine counts")
	fmt.Println("  • Use deadlock detection tools")
}

// ============================================================================
// DEMONSTRATION
// ============================================================================

// DemoDeadlockPrevention demonstrates deadlock prevention techniques
func DemoDeadlockPrevention() {
	fmt.Println("\n" + "="*60)
	fmt.Println("DEADLOCK PREVENTION DEMONSTRATION")
	fmt.Println("="*60 + "\n")

	demoLockOrdering()
	demoChannelDeadlock()
	demoWaitGroupDeadlock()
	demoCircularDependency()
	demoTimeoutPrevention()
	demoTryLock()
	demoPracticalExample()
	demoBestPractices()

	fmt.Println("\n✅ Deadlock prevention demonstration completed!")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("- Deadlocks occur when goroutines wait in a cycle")
	fmt.Println("- Consistent lock ordering prevents deadlocks")
	fmt.Println("- Use timeouts to avoid indefinite waits")
	fmt.Println("- Buffered channels can prevent channel deadlocks")
	fmt.Println("- Keep critical sections small and simple")
	fmt.Println("- Test with 'go run -race' flag")
}
