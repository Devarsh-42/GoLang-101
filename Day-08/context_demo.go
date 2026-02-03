package main

import (
	"context"
	"fmt"
	"time"
)

/*
CONTEXT PACKAGE
===============

What is it?
-----------
The context package provides a way to carry deadlines, cancellation signals,
and request-scoped values across API boundaries and between goroutines.
It's the standard way to manage goroutine lifecycles in Go.

Key Types:
- context.Context: Interface for carrying deadlines and cancellation
- context.Background(): Root context (never canceled)
- context.TODO(): Placeholder when context is unclear
- context.WithCancel(): Manual cancellation
- context.WithTimeout(): Time-based cancellation
- context.WithDeadline(): Absolute time cancellation
- context.WithValue(): Carry request-scoped data

Benefits:
- Prevents goroutine leaks
- Coordinates cancellation across multiple goroutines
- Implements timeouts consistently
- Passes request-scoped values safely

Use Cases:
- HTTP request handling (request timeout, cancellation)
- Database queries with timeout
- Graceful shutdown of services
- Canceling long-running operations
- Cascading cancellation in microservices
- Distributed tracing (passing trace IDs)

Best Practices:
1. Always pass context as first parameter
2. Don't store contexts in structs
3. Use context.Background() for main/init functions
4. Don't pass nil context (use context.TODO() if unsure)
5. Context values should be request-scoped data only
6. Cancel contexts to free resources
7. Check ctx.Done() in long-running operations
*/

// ============================================================================
// CONTEXT WITH CANCEL
// ============================================================================

// operationWithCancel demonstrates manual cancellation
func operationWithCancel(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("  Goroutine %d: Received cancellation signal: %v\n", id, ctx.Err())
			return
		default:
			fmt.Printf("  Goroutine %d: Working...\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func demoWithCancel() {
	fmt.Println("\n📌 Context with Cancel")
	fmt.Println("-" * 40)

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Start multiple goroutines
	for i := 1; i <= 3; i++ {
		go operationWithCancel(ctx, i)
	}

	// Let them run for 2 seconds
	time.Sleep(2 * time.Second)

	// Cancel all goroutines
	fmt.Println("\n🛑 Cancelling all goroutines...")
	cancel()

	// Give time for goroutines to receive cancellation
	time.Sleep(500 * time.Millisecond)
	fmt.Println("✓ All goroutines cancelled")
}

// ============================================================================
// CONTEXT WITH TIMEOUT
// ============================================================================

// slowOperation simulates a slow task
func slowOperation(ctx context.Context, duration time.Duration) error {
	done := make(chan bool)

	go func() {
		time.Sleep(duration)
		done <- true
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func demoWithTimeout() {
	fmt.Println("\n📌 Context with Timeout")
	fmt.Println("-" * 40)

	// Example 1: Operation completes before timeout
	fmt.Println("Test 1: Fast operation (500ms with 1s timeout)")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel1() // Always call cancel to release resources

	err := slowOperation(ctx1, 500*time.Millisecond)
	if err != nil {
		fmt.Println("  ✗ Error:", err)
	} else {
		fmt.Println("  ✓ Operation completed successfully")
	}

	// Example 2: Operation exceeds timeout
	fmt.Println("\nTest 2: Slow operation (2s with 1s timeout)")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel2()

	err = slowOperation(ctx2, 2*time.Second)
	if err != nil {
		fmt.Println("  ✗ Error:", err)
	} else {
		fmt.Println("  ✓ Operation completed successfully")
	}
}

// ============================================================================
// CONTEXT WITH DEADLINE
// ============================================================================

func demoWithDeadline() {
	fmt.Println("\n📌 Context with Deadline")
	fmt.Println("-" * 40)

	// Set deadline to 1 second from now
	deadline := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	fmt.Printf("Deadline set for: %s\n", deadline.Format("15:04:05"))
	fmt.Println("Starting operation...")

	err := slowOperation(ctx, 2*time.Second)
	if err != nil {
		fmt.Println("  ✗ Operation failed:", err)

		// Check if it's a deadline exceeded error
		if err == context.DeadlineExceeded {
			fmt.Println("  ℹ  Specifically: deadline exceeded")
		}
	} else {
		fmt.Println("  ✓ Operation completed")
	}
}

// ============================================================================
// CONTEXT WITH VALUE (Request-scoped data)
// ============================================================================

type contextKey string

const (
	userIDKey    contextKey = "userID"
	requestIDKey contextKey = "requestID"
)

// processRequest demonstrates passing values via context
func processRequest(ctx context.Context) {
	// Extract values from context
	userID := ctx.Value(userIDKey)
	requestID := ctx.Value(requestIDKey)

	fmt.Printf("  Processing request %v for user %v\n", requestID, userID)

	// Simulate passing context to another function
	authenticateUser(ctx)
	fetchData(ctx)
}

func authenticateUser(ctx context.Context) {
	userID := ctx.Value(userIDKey)
	fmt.Printf("  Authenticating user: %v\n", userID)
}

func fetchData(ctx context.Context) {
	requestID := ctx.Value(requestIDKey)
	fmt.Printf("  Fetching data for request: %v\n", requestID)
}

func demoWithValue() {
	fmt.Println("\n📌 Context with Value (Request-scoped data)")
	fmt.Println("-" * 40)

	// Create context with values
	ctx := context.Background()
	ctx = context.WithValue(ctx, userIDKey, "user123")
	ctx = context.WithValue(ctx, requestIDKey, "req456")

	fmt.Println("Simulating HTTP request processing:")
	processRequest(ctx)

	fmt.Println("\n⚠ Warning: Use context.WithValue sparingly!")
	fmt.Println("  Only for request-scoped data, not configuration")
}

// ============================================================================
// CONTEXT CHAINING (Timeout + Cancel)
// ============================================================================

func demoContextChaining() {
	fmt.Println("\n📌 Context Chaining")
	fmt.Println("-" * 40)

	// Parent context with timeout
	parentCtx, parentCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer parentCancel()

	// Child context with its own timeout
	childCtx, childCancel := context.WithTimeout(parentCtx, 1*time.Second)
	defer childCancel()

	fmt.Println("Parent timeout: 3s, Child timeout: 1s")
	fmt.Println("Starting operation with child context...")

	err := slowOperation(childCtx, 2*time.Second)
	if err != nil {
		fmt.Println("  ✗ Child context error:", err)
	}

	fmt.Println("✓ Child's timeout triggered first (1s < 3s)")
}

// ============================================================================
// PRACTICAL EXAMPLE: DATABASE QUERY WITH TIMEOUT
// ============================================================================

func simulateDBQuery(ctx context.Context, query string) error {
	fmt.Printf("  Executing query: %s\n", query)

	// Simulate query execution
	select {
	case <-time.After(800 * time.Millisecond):
		fmt.Println("  ✓ Query completed")
		return nil
	case <-ctx.Done():
		fmt.Println("  ✗ Query cancelled:", ctx.Err())
		return ctx.Err()
	}
}

func demoPracticalExample() {
	fmt.Println("\n📌 Practical Example: Database Query")
	fmt.Println("-" * 40)

	// Simulate HTTP request with 2-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Add request ID
	ctx = context.WithValue(ctx, requestIDKey, "http-req-789")

	requestID := ctx.Value(requestIDKey)
	fmt.Printf("Processing request: %v\n", requestID)

	// Execute multiple queries
	queries := []string{
		"SELECT * FROM users WHERE id = 1",
		"SELECT * FROM orders WHERE user_id = 1",
	}

	for _, query := range queries {
		if err := simulateDBQuery(ctx, query); err != nil {
			fmt.Printf("  Query failed: %v\n", err)
			break
		}
	}

	fmt.Println("✓ Request processing completed")
}

// ============================================================================
// DEMONSTRATION
// ============================================================================

// DemoContext demonstrates all context patterns
func DemoContext() {
	fmt.Println("\n" + "="*60)
	fmt.Println("CONTEXT PACKAGE DEMONSTRATION")
	fmt.Println("="*60 + "\n")

	demoWithCancel()
	demoWithTimeout()
	demoWithDeadline()
	demoWithValue()
	demoContextChaining()
	demoPracticalExample()

	fmt.Println("\n✅ Context demonstration completed!")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("- Context manages goroutine lifecycle and cancellation")
	fmt.Println("- WithTimeout/WithDeadline for time-based cancellation")
	fmt.Println("- WithCancel for manual cancellation")
	fmt.Println("- WithValue for request-scoped data (use sparingly)")
	fmt.Println("- Always call cancel() to release resources")
	fmt.Println("- Check ctx.Done() in long-running operations")
}
