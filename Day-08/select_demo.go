package main

import (
	"fmt"
	"time"
)

/*
SELECT STATEMENT
================

What is it?
-----------
The select statement lets a goroutine wait on multiple channel operations.
It blocks until one of its cases can proceed, then executes that case.
If multiple cases are ready, it chooses one at random.

Key Features:
- Waits on multiple channel operations simultaneously
- Blocks until at least one case can proceed
- Random selection when multiple cases are ready
- Default case provides non-blocking behavior
- Can be used in loops for continuous monitoring

Benefits:
- Multiplexing channel operations
- Implementing timeouts
- Non-blocking channel operations (with default)
- Graceful cancellation handling
- Coordination between multiple channels

Use Cases:
- Timeouts for operations
- Non-blocking sends/receives
- Coordinating multiple goroutines
- Implementing cancellation
- Monitoring multiple data sources
- Load balancing between channels

Best Practices:
1. Use default case for non-blocking operations
2. Combine with time.After() for timeouts
3. Use with context.Done() for cancellation
4. Be aware of random selection with multiple ready cases
5. Avoid blocking in case statements
*/

// ============================================================================
// BASIC SELECT EXAMPLES
// ============================================================================

// basicSelect demonstrates simple select usage
func basicSelect() {
	fmt.Println("\n📌 Basic Select Example")
	fmt.Println("-" * 40)

	ch1 := make(chan string)
	ch2 := make(chan string)

	// Send to ch1 after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "Message from channel 1"
	}()

	// Send to ch2 after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "Message from channel 2"
	}()

	// Select waits for the first channel to be ready
	select {
	case msg1 := <-ch1:
		fmt.Println("Received:", msg1)
	case msg2 := <-ch2:
		fmt.Println("Received:", msg2)
	}

	fmt.Println("✓ First channel responded")
}

// ============================================================================
// SELECT WITH DEFAULT (NON-BLOCKING)
// ============================================================================

// nonBlockingSelect demonstrates non-blocking channel operations
func nonBlockingSelect() {
	fmt.Println("\n📌 Non-Blocking Select (with default)")
	fmt.Println("-" * 40)

	messages := make(chan string)
	signals := make(chan bool)

	// Non-blocking receive
	select {
	case msg := <-messages:
		fmt.Println("Received message:", msg)
	default:
		fmt.Println("No message received (non-blocking)")
	}

	// Non-blocking send
	msg := "Hello"
	select {
	case messages <- msg:
		fmt.Println("Sent message:", msg)
	default:
		fmt.Println("No message sent (channel not ready)")
	}

	// Multiple non-blocking operations
	select {
	case msg := <-messages:
		fmt.Println("Received:", msg)
	case sig := <-signals:
		fmt.Println("Received signal:", sig)
	default:
		fmt.Println("No activity (non-blocking)")
	}

	fmt.Println("✓ Non-blocking operations completed")
}

// ============================================================================
// SELECT WITH TIMEOUT
// ============================================================================

// selectWithTimeout demonstrates timeout pattern
func selectWithTimeout() {
	fmt.Println("\n📌 Select with Timeout")
	fmt.Println("-" * 40)

	// Example 1: Operation completes before timeout
	ch1 := make(chan string, 1)
	go func() {
		time.Sleep(500 * time.Millisecond)
		ch1 <- "Fast operation completed"
	}()

	select {
	case result := <-ch1:
		fmt.Println("✓", result)
	case <-time.After(1 * time.Second):
		fmt.Println("✗ Operation timed out")
	}

	// Example 2: Operation exceeds timeout
	ch2 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "Slow operation completed"
	}()

	select {
	case result := <-ch2:
		fmt.Println("✓", result)
	case <-time.After(1 * time.Second):
		fmt.Println("✗ Operation timed out after 1 second")
	}
}

// ============================================================================
// SELECT IN LOOP (MULTIPLEXING)
// ============================================================================

// multiplexChannels demonstrates handling multiple channels continuously
func multiplexChannels() {
	fmt.Println("\n📌 Multiplexing Multiple Channels")
	fmt.Println("-" * 40)

	tick := time.Tick(500 * time.Millisecond)
	boom := time.After(2 * time.Second)

	fmt.Println("Starting ticker (stops after 2 seconds)...")

	for {
		select {
		case <-tick:
			fmt.Println("  tick.")
		case <-boom:
			fmt.Println("  BOOM!")
			fmt.Println("✓ Ticker stopped")
			return
		default:
			fmt.Print("    .")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// ============================================================================
// SELECT FOR QUIT/CANCELLATION CHANNEL
// ============================================================================

// quitChannelPattern demonstrates graceful shutdown
func quitChannelPattern() {
	fmt.Println("\n📌 Quit Channel Pattern")
	fmt.Println("-" * 40)

	quit := make(chan bool)
	data := make(chan int)

	// Producer goroutine
	go func() {
		for i := 1; i <= 10; i++ {
			select {
			case data <- i:
				fmt.Printf("Sent: %d\n", i)
				time.Sleep(200 * time.Millisecond)
			case <-quit:
				fmt.Println("Producer received quit signal")
				return
			}
		}
	}()

	// Consumer - receive a few items then quit
	for i := 0; i < 3; i++ {
		fmt.Printf("Received: %d\n", <-data)
	}

	// Send quit signal
	fmt.Println("\nSending quit signal...")
	quit <- true
	time.Sleep(100 * time.Millisecond)
	fmt.Println("✓ Graceful shutdown completed")
}

// ============================================================================
// SELECT WITH MULTIPLE READY CASES (RANDOM SELECTION)
// ============================================================================

// randomSelection demonstrates random case selection
func randomSelection() {
	fmt.Println("\n📌 Random Selection (Multiple Ready Cases)")
	fmt.Println("-" * 40)

	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	// Both channels are ready
	ch1 <- "Channel 1"
	ch2 <- "Channel 2"

	fmt.Println("Both channels have data ready. Running select 5 times:")
	for i := 1; i <= 5; i++ {
		// Re-fill channels
		select {
		case msg := <-ch1:
			fmt.Printf("  Iteration %d: Selected %s\n", i, msg)
			ch1 <- "Channel 1" // refill
		case msg := <-ch2:
			fmt.Printf("  Iteration %d: Selected %s\n", i, msg)
			ch2 <- "Channel 2" // refill
		}
	}

	fmt.Println("✓ Notice random selection when both cases are ready")
}

// ============================================================================
// DEMONSTRATION
// ============================================================================

// DemoSelect demonstrates all select patterns
func DemoSelect() {
	fmt.Println("\n" + "="*60)
	fmt.Println("SELECT STATEMENT DEMONSTRATION")
	fmt.Println("="*60 + "\n")

	basicSelect()
	nonBlockingSelect()
	selectWithTimeout()
	multiplexChannels()
	quitChannelPattern()
	randomSelection()

	fmt.Println("\n✅ Select demonstration completed!")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("- select blocks until one case can proceed")
	fmt.Println("- default case makes select non-blocking")
	fmt.Println("- time.After() enables timeout patterns")
	fmt.Println("- Random selection when multiple cases ready")
	fmt.Println("- Essential for coordinating multiple channels")
}
