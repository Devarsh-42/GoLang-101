# Day 07 - Concurrency in Go

## 📚 Topics Covered

Today's session focused on **Concurrency** in Go, one of the most powerful features of the language. We explored goroutines, channels, and concurrent programming patterns.

---

## 🔑 Key Concepts

### 1. **Goroutines**

Goroutines are lightweight threads managed by the Go runtime. They allow functions to run concurrently with other functions.

**Syntax:**
```go
go functionName()
```

**Example:**
```go
func hello() {
    fmt.Println("Hello, World!")
    time.Sleep(3 * time.Second)
    fmt.Println("Hello, again! Bbyee")
}

func main() {
    go hello() // Runs concurrently
    fmt.Println("Learning Concurrency in Go")
    time.Sleep(4 * time.Second)
}
```

**Key Points:**
- Goroutines are extremely cheap (only ~2KB of stack space initially)
- The main goroutine must wait for other goroutines to complete, or they'll be terminated when main exits
- Use `time.Sleep()` or channels for synchronization

![Goroutines Working](https://miro.medium.com/v2/resize:fit:1400/1*NFojvBkdPkoemSZ6binLbQ.png)

---

### 2. **Channels**

Channels are the pipes that connect concurrent goroutines. They allow goroutines to communicate and synchronize with each other.

**Types of Channels:**

#### **Unbuffered Channels**
- Blocking send and receive operations
- Sender blocks until receiver is ready
- Receiver blocks until sender sends data

```go
var requestChannel = make(chan string) // unbuffered channel

// Sending data
requestChannel <- "Hello from channel"

// Receiving data
msg := <-requestChannel
```

![Unbuffered Channel](https://miro.medium.com/v2/resize:fit:1400/1*gRHIPhVUiKGx2MQB0e7Kmg.png)

#### **Buffered Channels**
- Has a capacity to hold values
- Sender blocks only when buffer is full
- Receiver blocks only when buffer is empty

```go
var requestChannel3 = make(chan string, 1) // buffered channel with capacity 1
```

**Channel Operations:**
- **Send:** `channel <- value`
- **Receive:** `value := <-channel`
- **Close:** `close(channel)`

**Directional Channels:**
```go
// Send-only channel
func sendData(ch chan<- int) {
    ch <- 5
}

// Receive-only channel
func readData(ch <-chan int) {
    x := <-ch
    fmt.Println(x)
}
```

---

### 3. **Select Statement**

The `select` statement lets a goroutine wait on multiple channel operations. It blocks until one of its cases can proceed.

```go
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
```

**Key Points:**
- Executes the first case that is ready
- If multiple cases are ready, one is chosen randomly
- `default` case runs if no other case is ready (non-blocking)

---

### 4. **Done Channel Pattern**

A common pattern for synchronizing goroutine completion using a channel as a signaling mechanism.

```go
func doneChannelExample() {
    done := make(chan bool)

    go func() {
        time.Sleep(3 * time.Second)
        fmt.Println("Goroutine finished its work")
        done <- true // Signal completion
    }()

    fmt.Println("Waiting for goroutine to finish...")
    <-done // Wait for signal
    fmt.Println("Main function resumes after goroutine completion")
}
```

**Use Cases:**
- Waiting for background tasks to complete
- Coordinating multiple goroutines
- Graceful shutdown patterns

---

### 5. **Pipeline Pattern**

A pipeline is a series of stages connected by channels, where each stage is a group of goroutines running the same function.

```go
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
```

![Pipeline Working](https://miro.medium.com/v2/resize:fit:1400/1*8qBqyo8WW0gu6pzwCBwT_w.png)

**Key Points:**
- Each stage performs a specific transformation
- Data flows through channels from one stage to the next
- Close channels when no more data will be sent
- Use `range` to read from channels until closed

---

## 🎯 Code Examples Implemented

### 1. **Multiple Goroutines Running Concurrently**
```go
go hello()    // Prints messages with delays
go numbers()  // Prints numbers 1-5
go alphabets() // Prints letters a-e
```

### 2. **Channel Communication**
- Unbuffered channels for synchronization
- Buffered channels for asynchronous communication
- Closing channels to signal completion

### 3. **Select with Multiple Channels**
- Non-blocking channel operations
- Handling multiple channel sources

### 4. **Done Channel Synchronization**
- Proper goroutine coordination
- Waiting for task completion

### 5. **Pipeline for Data Processing**
- Three-stage pipeline
- Number generation → Squaring → Output

---

## 🚀 Running the Code

```bash
cd Day-07
go run main.go
```

---

## 📝 Key Takeaways

1. **Goroutines** enable lightweight concurrent execution
2. **Channels** provide safe communication between goroutines
3. **Unbuffered channels** synchronize goroutines (blocking)
4. **Buffered channels** allow asynchronous communication
5. **Select statements** handle multiple channel operations
6. **Done pattern** coordinates goroutine completion
7. **Pipelines** create efficient data processing flows
8. Always **close channels** when done sending to prevent deadlocks
9. The main goroutine should wait for child goroutines (using channels or sync primitives)

---

## 🔗 Resources

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go by Example - Goroutines](https://gobyexample.com/goroutines)
- [Go by Example - Channels](https://gobyexample.com/channels)

---

## 💡 Best Practices

1. **Don't communicate by sharing memory; share memory by communicating** (via channels)
2. Always close channels from the sender side, never the receiver
3. Use buffered channels when you know the capacity needed
4. Avoid goroutine leaks by ensuring they complete
5. Use `select` with `default` for non-blocking operations
6. Consider using `sync.WaitGroup` for managing multiple goroutines
7. Handle panics in goroutines to prevent crashes

---

## 🎓 Next Steps

- Explore `sync.WaitGroup` for better goroutine management
- Learn about `context` package for cancellation
- Study advanced patterns: Worker Pools, Fan-in/Fan-out
- Understand race conditions and the `race` detector
- Explore `sync.Mutex` for shared state protection
