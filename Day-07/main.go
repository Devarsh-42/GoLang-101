package main

import (
	"fmt"
	"time"
)

var requestChannel = make(chan string)     // unbuffered channel
var requestChannel2 = make(chan int)       // unbuffered channel of type int
var requestChannel3 = make(chan string, 1) // buffered channel - with capacity 1 because we have only one message to send

func hello() { // function to be run as a goroutine
	fmt.Println("Hello, World!")
	time.Sleep(3 * time.Second)
	fmt.Println("Hello, again! Bbyee")

	requestChannel <- "Hello function from channel"
	close(requestChannel)
}

// Example of send-only &  receive-only channel

// func sendData(ch chan<- int) {
//     ch <- 5// Sending data to the channel
// }

// func readData(ch <-chan int) {
//     x := <-ch // Receiving data from the channel
//     fmt.Println(x)
// }

// func main() {
//     ch := make(chan<- int) // Declaring a send-only channel

//     go sendData(ch) // Sending data to the channel
//     // go readData(ch) // This will cause a compile-time error

// ch2 := make(<-chan int) // Declaring a receive-only channel

// go readData(ch2) // Reading data from the channel

// }

func numbers() {
	for i := 1; i <= 5; i++ {
		requestChannel2 <- i // sending data to channel -> blocking operation
		time.Sleep(250 * time.Millisecond)
		fmt.Printf("%d ", i)
	}
	close(requestChannel2)
}
func alphabets() {
	for i := 'a'; i <= 'e'; i++ {
		time.Sleep(400 * time.Millisecond)
		fmt.Printf("%c ", i)
	}

	requestChannel3 <- "Alphabets function done"
	close(requestChannel3)
}

// DOne channel example
func doneChannelExample() {
	done := make(chan bool)

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("\nGoroutine finished its work")
		done <- true // signal that the goroutine is done
	}()

	fmt.Println("\nWaiting for goroutine to finish...")
	<-done // wait for the signal
	fmt.Println("\nMain function resumes after goroutine completion")
}

func main() {
	go hello() // goroutine
	fmt.Println("Learning Concurrency in Go")
	time.Sleep(4 * time.Second) // main goroutine waits for 5 seconds to let hello() finish
	fmt.Println("The Hello function has finished execution")

	go numbers()
	go alphabets()

	time.Sleep(3000 * time.Millisecond)
	fmt.Println("\nmain terminated without waiting for other goroutines")

	// Done channel example demonstration
	fmt.Println("\n--- Done Channel Example ---")
	doneChannelExample()

	// Using select with channels
	fmt.Println("\n--- Select Statement Example ---")

	select {
	case msg1 := <-requestChannel:
		fmt.Println("Received from hello goroutine:", msg1)
	case msg2 := <-requestChannel2:
		fmt.Println("Received from numbers goroutine:", msg2)
	case msg3 := <-requestChannel3:
		fmt.Println("Received from alphabets goroutine:", msg3)
	default:
		fmt.Println("No messages received")
	}


	// Pipline example
	fmt.Println("\n--- Pipeline Example ---")
	numbersChan := make(chan int)
	squaredChan := make(chan int)
	
	// Stage 1: Generate numbers
	go func() {
		for i := 1; i <= 5; i++ {
			numbersChan <- i
		}
		close(numbersChan)
	}()
	// Stage 2: Square the numbers
	go func() {
		for num := range numbersChan {
			squaredChan <- num * num
		}
		close(squaredChan)
	}()
	// Stage 3: Print the squared numbers
	for squaredNum := range squaredChan {
		fmt.Println("Squared Number:", squaredNum)
	}

	manav := &str{"Hello"}
	defer manav.p()
	manav.s = "World"
	fmt.Println(manav.s)
	panic("O no!")
}

type str struct {
	s string
}

func (s *str) p() {
	fmt.Println(s.s);
}