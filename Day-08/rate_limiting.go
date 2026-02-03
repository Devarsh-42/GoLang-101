package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
RATE LIMITING
=============

What is it?
-----------
Rate limiting controls the rate at which operations are executed.
It prevents resource exhaustion and ensures fair resource allocation.

Common Techniques:
1. Token Bucket: Tokens added at fixed rate, operations consume tokens
2. Leaky Bucket: Fixed rate of outflow
3. Fixed Window: Limit operations per time window
4. Sliding Window: More accurate than fixed window
5. Concurrency Limiting: Limit concurrent operations

Benefits:
- Prevents API rate limit violations
- Protects against resource exhaustion
- Ensures fair resource usage
- Prevents cascading failures
- Controls costs (for paid APIs)

Use Cases:
- API request throttling
- Database query rate limiting
- File processing rate control
- Preventing brute force attacks
- Controlling goroutine spawning
- Bandwidth management

Best Practices:
1. Choose appropriate algorithm for your use case
2. Make limits configurable
3. Provide feedback when limit is reached
4. Use context for cancellation
5. Consider burst capacity
6. Monitor and adjust limits based on metrics
*/

// ============================================================================
// TECHNIQUE 1: TIME.TICK (SIMPLE RATE LIMITER)
// ============================================================================

func demoSimpleRateLimiter() {
	fmt.Println("\n📌 Simple Rate Limiter (time.Tick)")
	fmt.Println("-" * 40)

	// Allow 1 request per 200ms (5 requests/second)
	limiter := time.Tick(200 * time.Millisecond)

	requests := []int{1, 2, 3, 4, 5}

	fmt.Println("Processing 5 requests at 5 req/sec...")
	start := time.Now()

	for _, req := range requests {
		<-limiter // Wait for rate limiter
		fmt.Printf("  Request %d processed at %v\n", req, time.Since(start).Round(time.Millisecond))
	}

	fmt.Printf("\n✓ All requests processed in %v\n", time.Since(start).Round(time.Millisecond))
	fmt.Println("  (Each request waited for rate limit)")
}

// ============================================================================
// TECHNIQUE 2: BUFFERED CHANNEL (TOKEN BUCKET)
// ============================================================================

type TokenBucket struct {
	capacity int
	tokens   chan struct{}
	fillRate time.Duration
	stopCh   chan struct{}
}

func NewTokenBucket(capacity int, fillRate time.Duration) *TokenBucket {
	tb := &TokenBucket{
		capacity: capacity,
		tokens:   make(chan struct{}, capacity),
		fillRate: fillRate,
		stopCh:   make(chan struct{}),
	}

	// Fill bucket initially
	for i := 0; i < capacity; i++ {
		tb.tokens <- struct{}{}
	}

	// Start token refill goroutine
	go tb.refill()

	return tb
}

func (tb *TokenBucket) refill() {
	ticker := time.NewTicker(tb.fillRate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case tb.tokens <- struct{}{}:
				// Token added
			default:
				// Bucket full, skip
			}
		case <-tb.stopCh:
			return
		}
	}
}

func (tb *TokenBucket) Take() bool {
	select {
	case <-tb.tokens:
		return true
	default:
		return false
	}
}

func (tb *TokenBucket) TakeWait(ctx context.Context) bool {
	select {
	case <-tb.tokens:
		return true
	case <-ctx.Done():
		return false
	}
}

func (tb *TokenBucket) Stop() {
	close(tb.stopCh)
}

func demoTokenBucket() {
	fmt.Println("\n📌 Token Bucket Rate Limiter")
	fmt.Println("-" * 40)

	// Capacity: 3 tokens, Refill: 1 token per 100ms
	bucket := NewTokenBucket(3, 100*time.Millisecond)
	defer bucket.Stop()

	fmt.Println("Bucket: 3 tokens, refill 1 token/100ms")
	fmt.Println("Sending 10 requests...\n")

	start := time.Now()

	for i := 1; i <= 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		if bucket.TakeWait(ctx) {
			elapsed := time.Since(start).Round(time.Millisecond)
			fmt.Printf("  ✓ Request %d processed at %v\n", i, elapsed)
		} else {
			fmt.Printf("  ✗ Request %d timed out\n", i)
		}

		cancel()
		time.Sleep(50 * time.Millisecond) // Simulate request interval
	}

	fmt.Printf("\n✓ Completed in %v\n", time.Since(start).Round(time.Millisecond))
	fmt.Println("  (Burst of 3, then rate-limited)")
}

// ============================================================================
// TECHNIQUE 3: CONCURRENCY LIMITER
// ============================================================================

type ConcurrencyLimiter struct {
	semaphore chan struct{}
}

func NewConcurrencyLimiter(limit int) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{
		semaphore: make(chan struct{}, limit),
	}
}

func (cl *ConcurrencyLimiter) Acquire() {
	cl.semaphore <- struct{}{}
}

func (cl *ConcurrencyLimiter) Release() {
	<-cl.semaphore
}

func (cl *ConcurrencyLimiter) Execute(fn func()) {
	cl.Acquire()
	defer cl.Release()
	fn()
}

func demoConcurrencyLimiter() {
	fmt.Println("\n📌 Concurrency Limiter")
	fmt.Println("-" * 40)

	limiter := NewConcurrencyLimiter(3) // Max 3 concurrent operations
	var wg sync.WaitGroup

	fmt.Println("Limiting concurrent operations to 3...")
	start := time.Now()

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			limiter.Execute(func() {
				elapsed := time.Since(start).Round(time.Millisecond)
				fmt.Printf("  [%v] Task %d started (concurrent)\n", elapsed, id)
				time.Sleep(200 * time.Millisecond)
				fmt.Printf("  [%v] Task %d finished\n", time.Since(start).Round(time.Millisecond), id)
			})
		}(i)
	}

	wg.Wait()
	fmt.Printf("\n✓ All tasks completed in %v\n", time.Since(start).Round(time.Millisecond))
	fmt.Println("  (Max 3 concurrent tasks at any time)")
}

// ============================================================================
// TECHNIQUE 4: RATE LIMITER WITH GOLANG.ORG/X/TIME/RATE
// (Simulated implementation)
// ============================================================================

type RateLimiter struct {
	rate       time.Duration
	burst      int
	tokens     int
	lastRefill time.Time
	mu         sync.Mutex
}

func NewRateLimiter(requestsPerSecond int, burst int) *RateLimiter {
	return &RateLimiter{
		rate:       time.Second / time.Duration(requestsPerSecond),
		burst:      burst,
		tokens:     burst,
		lastRefill: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	newTokens := int(elapsed / rl.rate)

	if newTokens > 0 {
		rl.tokens += newTokens
		if rl.tokens > rl.burst {
			rl.tokens = rl.burst
		}
		rl.lastRefill = now
	}

	// Check if we have tokens
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

func demoRateLimiter() {
	fmt.Println("\n📌 Advanced Rate Limiter (Token Bucket)")
	fmt.Println("-" * 40)

	// 5 requests per second, burst of 10
	limiter := NewRateLimiter(5, 10)

	fmt.Println("Rate: 5 req/sec, Burst: 10")
	fmt.Println("Sending 20 requests rapidly...\n")

	allowed := 0
	rejected := 0
	start := time.Now()

	for i := 1; i <= 20; i++ {
		if limiter.Allow() {
			allowed++
			fmt.Printf("  ✓ Request %d allowed at %v\n", i, time.Since(start).Round(time.Millisecond))
		} else {
			rejected++
			fmt.Printf("  ✗ Request %d rejected at %v\n", i, time.Since(start).Round(time.Millisecond))
		}
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Printf("\nResults: %d allowed, %d rejected\n", allowed, rejected)
	fmt.Println("✓ Burst allows initial requests, then rate limits")
}

// ============================================================================
// PRACTICAL EXAMPLE: API CLIENT WITH RATE LIMITING
// ============================================================================

type APIClient struct {
	limiter *ConcurrencyLimiter
	bucket  *TokenBucket
}

func NewAPIClient(maxConcurrent int, requestsPerSecond int) *APIClient {
	return &APIClient{
		limiter: NewConcurrencyLimiter(maxConcurrent),
		bucket:  NewTokenBucket(requestsPerSecond, time.Second/time.Duration(requestsPerSecond)),
	}
}

func (c *APIClient) MakeRequest(ctx context.Context, id int) error {
	// Wait for rate limit
	if !c.bucket.TakeWait(ctx) {
		return fmt.Errorf("request %d: rate limit timeout", id)
	}

	// Acquire concurrency slot
	c.limiter.Acquire()
	defer c.limiter.Release()

	// Simulate API call
	fmt.Printf("  📡 Request %d: Making API call...\n", id)
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("  ✓ Request %d: Success\n", id)

	return nil
}

func (c *APIClient) Close() {
	c.bucket.Stop()
}

func demoPracticalExample() {
	fmt.Println("\n📌 Practical Example: API Client")
	fmt.Println("-" * 40)

	// Max 3 concurrent requests, 5 requests per second
	client := NewAPIClient(3, 5)
	defer client.Close()

	fmt.Println("API Client: 3 concurrent, 5 req/sec")
	fmt.Println("Sending 10 requests...\n")

	var wg sync.WaitGroup
	start := time.Now()

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := client.MakeRequest(ctx, id); err != nil {
				fmt.Printf("  ✗ %v\n", err)
			}
		}(i)
	}

	wg.Wait()

	fmt.Printf("\n✓ All requests completed in %v\n", time.Since(start).Round(time.Millisecond))
	fmt.Println("  (Respects both rate limit and concurrency limit)")
}

// ============================================================================
// BEST PRACTICES SUMMARY
// ============================================================================

func demoBestPractices() {
	fmt.Println("\n📌 Rate Limiting Best Practices")
	fmt.Println("-" * 40)

	fmt.Println("\n✅ Choose the right technique:")
	fmt.Println("  • time.Tick: Simple, fixed rate")
	fmt.Println("  • Token Bucket: Allows bursts, smooth rate")
	fmt.Println("  • Concurrency Limiter: Limit concurrent ops")
	fmt.Println("  • Combined: Both rate + concurrency limits")

	fmt.Println("\n✅ Implementation tips:")
	fmt.Println("  • Make limits configurable")
	fmt.Println("  • Use context for timeouts")
	fmt.Println("  • Provide clear error messages")
	fmt.Println("  • Monitor rate limit hits")
	fmt.Println("  • Consider exponential backoff")

	fmt.Println("\n✅ Common use cases:")
	fmt.Println("  • API clients: Respect provider limits")
	fmt.Println("  • Servers: Protect from overload")
	fmt.Println("  • Workers: Control resource usage")
	fmt.Println("  • Database: Prevent connection exhaustion")
}

// ============================================================================
// DEMONSTRATION
// ============================================================================

// DemoRateLimiting demonstrates all rate limiting patterns
func DemoRateLimiting() {
	fmt.Println("\n" + "="*60)
	fmt.Println("RATE LIMITING DEMONSTRATION")
	fmt.Println("="*60 + "\n")

	demoSimpleRateLimiter()
	demoTokenBucket()
	demoConcurrencyLimiter()
	demoRateLimiter()
	demoPracticalExample()
	demoBestPractices()

	fmt.Println("\n✅ Rate limiting demonstration completed!")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("- Rate limiting prevents resource exhaustion")
	fmt.Println("- Token bucket allows controlled bursts")
	fmt.Println("- Concurrency limiter controls parallel operations")
	fmt.Println("- Combine techniques for robust rate limiting")
	fmt.Println("- Always use context for cancellation")
}
