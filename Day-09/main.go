package main

import (
	"context"
	"fmt"
	"time"
)

// type Site struct {
// 	URL string
// }


// type Response struct {
// 	Status int
// }

func main() {
	ctx := context.Background()

	exampleTimeout(ctx)

}

func exampleTimeout(ctx context.Context) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	done := make(chan struct{}) // channel to signal completion of work 

	go func() { // Simulate some work that takes time
		time.Sleep(2 * time.Second)
		fmt.Println("Finished work")
		close(done)
	}()	

	select {
		case <-done: // done channel is closed when work is completed
			fmt.Println("Work completed within timeout")
		case <-ctxWithTimeout.Done(): // ctxWithTimeout.Done() is closed when timeout is exceeded
			fmt.Println("Timeout exceeded:", ctxWithTimeout.Err())
	}
}