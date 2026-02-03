package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

/*
GRACEFUL SHUTDOWN
=================

What is it?
-----------
Graceful shutdown is the process of cleanly stopping a running application,
ensuring that:
1. No new work is accepted
2. In-flight work completes or times out
3. Resources are properly released
4. State is saved if needed
5. Connections are closed cleanly

Why is it important?
- Prevents data loss
- Avoids corrupted state
- Proper resource cleanup
- Better user experience
- No orphaned goroutines

Common Scenarios:
- HTTP servers handling active requests
- Message queue consumers processing messages
- Database connections and transactions
- File operations and uploads
- Long-running batch jobs

Implementation Steps:
1. Listen for shutdown signals (SIGINT, SIGTERM)
2. Stop accepting new work
3. Wait for in-flight work with timeout
4. Clean up resources (close connections, files, etc.)
5. Exit gracefully

Best Practices:
1. Use context for cancellation propagation
2. Set reasonable shutdown timeouts
3. Use WaitGroup to track in-flight work
4. Log shutdown process for debugging
5. Handle signals properly (SIGTERM, SIGINT)
6. Save critical state before exit
7. Return appropriate exit codes
*/

// ============================================================================
// EXAMPLE 1: BASIC GRACEFUL SHUTDOWN
// ============================================================================

func demoBasicGracefulShutdown() {
	fmt.Println("\n📌 Basic Graceful Shutdown")
	fmt.Println("-" * 40)

	// Create shutdown channel
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Simulate work
	done := make(chan bool)

	go func() {
		fmt.Println("Worker: Starting...")

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("  Worker: Processing...")
			case <-shutdown:
				fmt.Println("  Worker: Received shutdown signal")
				fmt.Println("  Worker: Cleaning up...")
				time.Sleep(200 * time.Millisecond)
				fmt.Println("  Worker: Finished cleanup")
				done <- true
				return
			}
		}
	}()

	fmt.Println("\nSimulating work for 2 seconds, then shutting down...")
	time.Sleep(2 * time.Second)

	// Trigger shutdown
	shutdown <- syscall.SIGTERM

	// Wait for graceful shutdown
	<-done
	fmt.Println("✓ Graceful shutdown completed")
}

// ============================================================================
// EXAMPLE 2: SHUTDOWN WITH CONTEXT
// ============================================================================

func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d: Started\n", id)

	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("  Worker %d: Shutdown signal received\n", id)
			// Simulate cleanup
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("  Worker %d: Cleanup completed\n", id)
			return
		case <-ticker.C:
			fmt.Printf("  Worker %d: Processing...\n", id)
		}
	}
}

func demoContextShutdown() {
	fmt.Println("\n📌 Graceful Shutdown with Context")
	fmt.Println("-" * 40)

	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Start multiple workers
	numWorkers := 3
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, &wg)
	}

	// Simulate running for a while
	fmt.Printf("\nStarting %d workers...\n", numWorkers)
	time.Sleep(1500 * time.Millisecond)

	// Trigger shutdown
	fmt.Println("\n🛑 Initiating graceful shutdown...")
	cancel() // Cancel context - all workers receive signal

	// Wait for all workers to finish
	wg.Wait()

	fmt.Println("✓ All workers shut down gracefully")
}

// ============================================================================
// EXAMPLE 3: SHUTDOWN WITH TIMEOUT
// ============================================================================

func taskWithTimeout(ctx context.Context, id int, duration time.Duration, results chan<- string) {
	taskDone := make(chan bool)

	go func() {
		fmt.Printf("  Task %d: Processing (will take %v)...\n", id, duration)
		time.Sleep(duration)
		taskDone <- true
	}()

	select {
	case <-taskDone:
		results <- fmt.Sprintf("Task %d completed successfully", id)
	case <-ctx.Done():
		results <- fmt.Sprintf("Task %d cancelled (timeout/shutdown)", id)
	}
}

func demoShutdownWithTimeout() {
	fmt.Println("\n📌 Graceful Shutdown with Timeout")
	fmt.Println("-" * 40)

	// Context with 2-second shutdown timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	results := make(chan string, 5)

	// Start tasks with different durations
	fmt.Println("Starting tasks with various durations...")
	go taskWithTimeout(ctx, 1, 500*time.Millisecond, results)  // Will complete
	go taskWithTimeout(ctx, 2, 1*time.Second, results)         // Will complete
	go taskWithTimeout(ctx, 3, 1500*time.Millisecond, results) // Will complete
	go taskWithTimeout(ctx, 4, 3*time.Second, results)         // Will timeout
	go taskWithTimeout(ctx, 5, 5*time.Second, results)         // Will timeout

	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n🛑 Shutdown initiated with 2-second timeout...")

	// Collect results
	for i := 0; i < 5; i++ {
		result := <-results
		fmt.Printf("  %s\n", result)
	}

	fmt.Println("\n✓ Shutdown completed (some tasks cancelled due to timeout)")
}

// ============================================================================
// EXAMPLE 4: SERVER-LIKE GRACEFUL SHUTDOWN
// ============================================================================

type Server struct {
	shutdown chan struct{}
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewServer() *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		shutdown: make(chan struct{}),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (s *Server) handleRequest(id int) {
	defer s.wg.Done()

	fmt.Printf("  🌐 Handling request %d...\n", id)

	// Simulate request processing
	processDone := make(chan bool)
	go func() {
		time.Sleep(time.Duration(id%3+1) * 200 * time.Millisecond)
		processDone <- true
	}()

	select {
	case <-processDone:
		fmt.Printf("  ✓ Request %d completed\n", id)
	case <-s.ctx.Done():
		fmt.Printf("  ⚠ Request %d cancelled (shutdown)\n", id)
	}
}

func (s *Server) Start() {
	fmt.Println("Server: Starting...")

	// Simulate accepting requests
	go func() {
		requestID := 1
		ticker := time.NewTicker(300 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.wg.Add(1)
				go s.handleRequest(requestID)
				requestID++
			case <-s.shutdown:
				fmt.Println("\nServer: Stopped accepting new requests")
				return
			}
		}
	}()
}

func (s *Server) Shutdown(timeout time.Duration) error {
	fmt.Printf("\n🛑 Server: Initiating graceful shutdown (timeout: %v)...\n", timeout)

	// Stop accepting new requests
	close(s.shutdown)

	// Cancel context for in-flight requests
	s.cancel()

	// Wait for in-flight requests with timeout
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("Server: All requests completed")
		return nil
	case <-time.After(timeout):
		fmt.Println("Server: Shutdown timeout - forcing exit")
		return fmt.Errorf("shutdown timeout exceeded")
	}
}

func demoServerShutdown() {
	fmt.Println("\n📌 Server-Like Graceful Shutdown")
	fmt.Println("-" * 40)

	server := NewServer()
	server.Start()

	// Let server handle some requests
	time.Sleep(2 * time.Second)

	// Graceful shutdown
	if err := server.Shutdown(3 * time.Second); err != nil {
		fmt.Printf("✗ Shutdown error: %v\n", err)
	} else {
		fmt.Println("✓ Server shut down gracefully")
	}
}

// ============================================================================
// EXAMPLE 5: MULTI-COMPONENT SHUTDOWN
// ============================================================================

type Component struct {
	name     string
	shutdown chan struct{}
	done     chan struct{}
}

func NewComponent(name string) *Component {
	c := &Component{
		name:     name,
		shutdown: make(chan struct{}),
		done:     make(chan struct{}),
	}
	go c.run()
	return c
}

func (c *Component) run() {
	defer close(c.done)

	fmt.Printf("  %s: Started\n", c.name)
	ticker := time.NewTicker(400 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("  %s: Working...\n", c.name)
		case <-c.shutdown:
			fmt.Printf("  %s: Shutting down...\n", c.name)
			time.Sleep(200 * time.Millisecond) // Cleanup
			fmt.Printf("  %s: Shutdown complete\n", c.name)
			return
		}
	}
}

func (c *Component) Shutdown() {
	close(c.shutdown)
	<-c.done
}

func demoMultiComponentShutdown() {
	fmt.Println("\n📌 Multi-Component Graceful Shutdown")
	fmt.Println("-" * 40)

	// Start multiple components
	components := []*Component{
		NewComponent("Database"),
		NewComponent("Cache"),
		NewComponent("MessageQueue"),
		NewComponent("HTTPServer"),
	}

	fmt.Println("\nAll components running...")
	time.Sleep(1500 * time.Millisecond)

	// Shutdown all components
	fmt.Println("\n🛑 Shutting down all components...")

	var wg sync.WaitGroup
	for _, comp := range components {
		wg.Add(1)
		go func(c *Component) {
			defer wg.Done()
			c.Shutdown()
		}(comp)
	}

	wg.Wait()
	fmt.Println("\n✓ All components shut down gracefully")
}

// ============================================================================
// BEST PRACTICES
// ============================================================================

func demoBestPractices() {
	fmt.Println("\n📌 Graceful Shutdown Best Practices")
	fmt.Println("-" * 40)

	fmt.Println("\n✅ Signal Handling:")
	fmt.Println("  • Listen for SIGINT and SIGTERM")
	fmt.Println("  • Use buffered channel for signals")
	fmt.Println("  • Don't ignore SIGTERM (used by orchestrators)")

	fmt.Println("\n✅ Context Usage:")
	fmt.Println("  • Use context for cancellation propagation")
	fmt.Println("  • Set reasonable timeouts")
	fmt.Println("  • Check ctx.Done() in long operations")

	fmt.Println("\n✅ Resource Cleanup:")
	fmt.Println("  • Close database connections")
	fmt.Println("  • Flush buffers and caches")
	fmt.Println("  • Save critical state")
	fmt.Println("  • Close file handles")

	fmt.Println("\n✅ Shutdown Sequence:")
	fmt.Println("  1. Stop accepting new work")
	fmt.Println("  2. Signal existing workers to stop")
	fmt.Println("  3. Wait for workers (with timeout)")
	fmt.Println("  4. Clean up resources")
	fmt.Println("  5. Exit with appropriate code")

	fmt.Println("\n✅ Testing:")
	fmt.Println("  • Test shutdown under load")
	fmt.Println("  • Verify no data loss")
	fmt.Println("  • Check for goroutine leaks")
	fmt.Println("  • Test timeout scenarios")
}

// ============================================================================
// DEMONSTRATION
// ============================================================================

// DemoGracefulShutdown demonstrates graceful shutdown patterns
func DemoGracefulShutdown() {
	fmt.Println("\n" + "="*60)
	fmt.Println("GRACEFUL SHUTDOWN DEMONSTRATION")
	fmt.Println("="*60 + "\n")

	demoBasicGracefulShutdown()
	demoContextShutdown()
	demoShutdownWithTimeout()
	demoServerShutdown()
	demoMultiComponentShutdown()
	demoBestPractices()

	fmt.Println("\n✅ Graceful shutdown demonstration completed!")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("- Graceful shutdown prevents data loss")
	fmt.Println("- Use context for cancellation propagation")
	fmt.Println("- Set reasonable shutdown timeouts")
	fmt.Println("- Wait for in-flight work to complete")
	fmt.Println("- Clean up resources properly")
	fmt.Println("- Handle SIGINT and SIGTERM signals")
}
