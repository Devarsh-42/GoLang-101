package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var wg sync.WaitGroup
var mu sync.Mutex

func SendRequest(url string) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching %s: %v\n", url, err)
		panic(err)
	}
	mu.Lock()
	fmt.Printf("[%d] Response from %s\n", resp.StatusCode, url)
	defer mu.Unlock()
	defer resp.Body.Close()

}

func main() {

	if len(os.Args) < 2 {
		log.Fatalln("Please provide a filename as an argument")
		return
	}

	for _, url := range os.Args[1:] {
		go SendRequest("http://" + url)
		wg.Add(1)
	}
	wg.Wait()

}
