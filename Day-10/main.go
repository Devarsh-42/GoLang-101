package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// =============================================================================
// Advanced HTTP Client - Common patterns used in real development
// =============================================================================

// RequestResult holds the response data from an HTTP request
type RequestResult struct {
	URL        string
	StatusCode int
	BodySize   int
	Duration   time.Duration
	Error      error
}

var (
	wg sync.WaitGroup
	mu sync.Mutex
)

// SendRequest makes an HTTP request with context support (for timeout/cancellation)
func SendRequest(ctx context.Context, url string) RequestResult {
	start := time.Now()

	// Create a request with context - this allows timeout & cancellation
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return RequestResult{URL: url, Error: err, Duration: time.Since(start)}
	}

	// Use a custom HTTP client with timeout (don't use http.Get in production!)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return RequestResult{URL: url, Error: err, Duration: time.Since(start)}
	}
	defer resp.Body.Close()

	// Read the body to get size
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RequestResult{URL: url, StatusCode: resp.StatusCode, Error: err, Duration: time.Since(start)}
	}

	return RequestResult{
		URL:        url,
		StatusCode: resp.StatusCode,
		BodySize:   len(body),
		Duration:   time.Since(start),
	}
}

// SendRequestWithRetry retries a failed request up to maxRetries times
// This is a very common pattern in production code
func SendRequestWithRetry(ctx context.Context, url string, maxRetries int) RequestResult {
	var result RequestResult

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: wait longer between each retry
			// attempt 1 = 1s, attempt 2 = 2s, attempt 3 = 4s
			backoff := time.Duration(1<<(attempt-1)) * time.Second
			log.Printf("⏳ Retry #%d for %s (waiting %v)\n", attempt, url, backoff)

			select {
			case <-ctx.Done():
				result.Error = ctx.Err()
				return result
			case <-time.After(backoff):
			}
		}

		result = SendRequest(ctx, url)
		if result.Error == nil && result.StatusCode < 500 {
			return result // success or client error (4xx) - don't retry
		}
	}

	return result
}

// FetchAllConcurrent fetches multiple URLs concurrently with a worker limit
// This prevents overwhelming the server with too many connections
func FetchAllConcurrent(ctx context.Context, urls []string, maxConcurrent int) []RequestResult {
	results := make([]RequestResult, len(urls))

	// Semaphore pattern: limits how many goroutines run at once
	semaphore := make(chan struct{}, maxConcurrent)

	for i, url := range urls {
		wg.Add(1)
		go func(index int, u string) {
			defer wg.Done()

			semaphore <- struct{}{}        // acquire slot (blocks if full)
			defer func() { <-semaphore }() // release slot when done

			result := SendRequestWithRetry(ctx, u, 2) // retry up to 2 times

			mu.Lock()
			results[index] = result
			mu.Unlock()
		}(i, url)
	}

	wg.Wait()
	return results
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: go run main.go <url1> <url2> ...")
	}

	// Prepare URLs
	urls := make([]string, len(os.Args[1:]))
	for i, arg := range os.Args[1:] {
		urls[i] = "http://" + arg
	}

	// Create context with overall timeout of 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("🚀 Fetching URLs concurrently (max 3 at a time)...")
	start := time.Now()

	// Fetch all URLs with max 3 concurrent requests
	results := FetchAllConcurrent(ctx, urls, 3)

	// Print results
	fmt.Println("📊 Results:")
	fmt.Println("─────────────────────────────────────────────────")
	successCount := 0
	failCount := 0
	for _, r := range results {
		if r.Error != nil {
			fmt.Printf("  ❌ %s → Error: %v (%v)\n", r.URL, r.Error, r.Duration)
			failCount++
		} else {
			fmt.Printf("  ✅ [%d] %s → %d bytes (%v)\n", r.StatusCode, r.URL, r.BodySize, r.Duration)
			successCount++
		}
	}
	fmt.Println("─────────────────────────────────────────────────")
	fmt.Printf("  Total: %d URLs | ✅ %d success | ❌ %d failed | ⏱️  %v\n",
		len(results), successCount, failCount, time.Since(start))
}
