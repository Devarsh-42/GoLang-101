package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// =============================================================================
// TESTS - Verify correctness (does it work?)
// =============================================================================

// TestSendRequestSuccess tests a basic successful HTTP request
func TestSendRequestSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from test server"))
	}))
	defer server.Close()

	ctx := context.Background()
	result := SendRequest(ctx, server.URL)

	if result.Error != nil {
		t.Fatalf("Expected no error, got: %v", result.Error)
	}
	if result.StatusCode != 200 {
		t.Fatalf("Expected status 200, got: %d", result.StatusCode)
	}
	if result.BodySize != 22 { // "Hello from test server" = 22 bytes
		t.Fatalf("Expected body size 22, got: %d", result.BodySize)
	}
	t.Logf("✅ Request succeeded: %d bytes in %v", result.BodySize, result.Duration)
}

// TestSendRequestTimeout tests that context timeout cancels the request
func TestSendRequestTimeout(t *testing.T) {
	// Server that takes 5 seconds to respond
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.Write([]byte("too slow"))
	}))
	defer server.Close()

	// Context with 100ms timeout - request should fail
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result := SendRequest(ctx, server.URL)

	if result.Error == nil {
		t.Fatal("Expected timeout error, got nil")
	}
	t.Logf("✅ Request correctly timed out: %v", result.Error)
}

// TestSendRequestBadURL tests handling of an invalid URL
func TestSendRequestBadURL(t *testing.T) {
	ctx := context.Background()
	result := SendRequest(ctx, "http://localhost:99999/bad")

	if result.Error == nil {
		t.Fatal("Expected error for bad URL, got nil")
	}
	t.Logf("✅ Bad URL correctly returned error: %v", result.Error)
}

// TestSendRequestWithRetrySuccess tests retry succeeds on second attempt
func TestSendRequestWithRetrySuccess(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 1 {
			w.WriteHeader(http.StatusInternalServerError) // 500 on first try
			return
		}
		w.WriteHeader(http.StatusOK) // 200 on retry
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	ctx := context.Background()
	result := SendRequestWithRetry(ctx, server.URL, 2)

	if result.Error != nil {
		t.Fatalf("Expected no error after retry, got: %v", result.Error)
	}
	if result.StatusCode != 200 {
		t.Fatalf("Expected status 200 after retry, got: %d", result.StatusCode)
	}
	t.Logf("✅ Retry worked! Took %d attempts", attempts)
}

// TestFetchAllConcurrent tests fetching multiple URLs at once
func TestFetchAllConcurrent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("response"))
	}))
	defer server.Close()

	urls := []string{server.URL, server.URL, server.URL}
	ctx := context.Background()

	results := FetchAllConcurrent(ctx, urls, 2)

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got: %d", len(results))
	}
	for i, r := range results {
		if r.Error != nil {
			t.Fatalf("Request %d failed: %v", i, r.Error)
		}
		if r.StatusCode != 200 {
			t.Fatalf("Request %d: expected 200, got %d", i, r.StatusCode)
		}
	}
	t.Logf("✅ All %d concurrent requests succeeded", len(results))
}

// TestFetchAllConcurrentWithTimeout tests that timeout cancels all requests
func TestFetchAllConcurrentWithTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second) // very slow server
		w.Write([]byte("slow"))
	}))
	defer server.Close()

	urls := []string{server.URL, server.URL}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	results := FetchAllConcurrent(ctx, urls, 2)

	for _, r := range results {
		if r.Error == nil {
			t.Fatal("Expected timeout error, got nil")
		}
	}
	t.Log("✅ All requests correctly timed out")
}

// =============================================================================
// BENCHMARKS - Measure performance (how fast is it?)
// Run with: go test -bench=. -benchmem
// =============================================================================

// BenchmarkSendRequest benchmarks a single HTTP request
func BenchmarkSendRequest(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("benchmark response"))
	}))
	defer server.Close()

	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := SendRequest(ctx, server.URL)
		if result.Error != nil {
			b.Fatalf("Error: %v", result.Error)
		}
	}
}

// BenchmarkFetchAllConcurrent benchmarks concurrent fetching
func BenchmarkFetchAllConcurrent(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("benchmark response"))
	}))
	defer server.Close()

	urls := []string{server.URL, server.URL, server.URL, server.URL, server.URL}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FetchAllConcurrent(ctx, urls, 3)
	}
}
