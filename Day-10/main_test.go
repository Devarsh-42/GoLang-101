package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkSendRequest benchmarks a single HTTP request
func BenchmarkSendRequest(b *testing.B) {
	// Create a test server that responds immediately
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Remove the panic to prevent benchmark from stopping
		resp, err := http.Get(server.URL)
		if err != nil {
			b.Fatalf("Error fetching: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkSendRequestParallel benchmarks parallel HTTP requests
func BenchmarkSendRequestParallel(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := http.Get(server.URL)
			if err != nil {
				b.Fatalf("Error fetching: %v", err)
			}
			resp.Body.Close()
		}
	})
}

// BenchmarkMultipleRequests benchmarks multiple sequential requests
func BenchmarkMultipleRequests(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	urls := []string{server.URL, server.URL, server.URL, server.URL, server.URL}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, url := range urls {
			resp, err := http.Get(url)
			if err != nil {
				b.Fatalf("Error fetching: %v", err)
			}
			resp.Body.Close()
		}
	}
}
