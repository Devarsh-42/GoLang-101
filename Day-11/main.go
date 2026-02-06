package main

import (
	"log/slog"
	"net/http"
)


func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}",handleHello)
	mux.HandleFunc("/goodbye",handleGoodBye)
	

	http.ListenAndServe(":8080", mux) 

}

func handleHello(w http.ResponseWriter, r *http.Request) {
	wc,err := w.Write([]byte("Hello, World!"))
	if err != nil {
		slog.Error("Failed to write response", "error", err)
		return
	}
	println("Bytes written:", wc)

}

func handleGoodBye(w http.ResponseWriter, r *http.Request) {
	wc,err := w.Write([]byte("Goodbye, World!")) // Write the response body and check for errors, wc = number of bytes written
	if err != nil {
		slog.Error("Failed to write response", "error", err)
		return
	}
	println("Bytes written:", wc)
}