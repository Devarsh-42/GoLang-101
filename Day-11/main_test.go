package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHello(t *testing.T) {
	// Create a request to pass to our handler
	req := httptest.NewRequest("GET", "http://localhost:8080/", nil)
	w := httptest.NewRecorder()

	// Call the handler function directly
	handleHello(w, req)

	// Check the response
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK but got %v", resp.Status)
	}
	if string(body) != "Hello, World!" {
		t.Errorf("Expected body 'Hello, World!' but got '%s'", string(body))
	}
}


func TestHandleGoodBye(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/goodbye", nil)
	w := httptest.NewRecorder()

	handleGoodBye(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK but got %v", resp.Status)
	}
	if string(body) != "Goodbye, World!" {
		t.Errorf("Expected body 'Goodbye, World!' but got '%s'", string(body))
	}
}