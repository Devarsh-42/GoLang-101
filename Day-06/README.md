# Day 06 - Go Concurrency

## Topics Covered

### 1. Concurrency vs Parallelism
- **Concurrency**: Dealing with multiple tasks at once (composition of independently executing processes)
- **Parallelism**: Doing multiple tasks simultaneously (requires multiple CPU cores)
- Go is designed for concurrency, parallelism is a consequence when run on multiple cores

### 2. Goroutines
- Lightweight threads managed by Go runtime
- Created using `go` keyword: `go functionName()`
- Much cheaper than OS threads (start with ~2KB stack)
- Thousands of goroutines can run simultaneously
- Example:
  ```go
  go func() {
      // code runs concurrently
  }()
  ```

### 3. Go Concurrency Structure
- **CSP (Communicating Sequential Processes)**: Go's concurrency model
- "Don't communicate by sharing memory; share memory by communicating"
- Goroutines communicate via channels, not shared variables

### 4. Channels
- Typed conduits for communication between goroutines
- Created with `make()`: `ch := make(chan int)`
- Send: `ch <- value`
- Receive: `value := <-ch`
- Blocking by default (synchronous)
- Close channels: `close(ch)`

### 5. Buffered Channels
- Channels with capacity: `ch := make(chan int, 100)`
- Non-blocking sends until buffer is full
- Non-blocking receives until buffer is empty
- Useful for throttling and batching

### 6. Select Statement
- Multiplexes multiple channel operations
- Waits on multiple channels simultaneously
- Executes whichever case is ready first
- `default` case prevents blocking
- Example:
  ```go
  select {
  case v := <-ch1:
      // handle ch1
  case ch2 <- value:
      // send to ch2
  default:
      // non-blocking
  }
  ```

### 7. For-Select Pattern
- Common idiom for continuous channel processing
- Combines `for` loop with `select`
- Example:
  ```go
  for {
      select {
      case v := <-ch:
          // process value
      case <-done:
          return
      }
  }
  ```

### 8. Done Channel & Pipeline
- **Done Channel**: Signal goroutine termination
  - Closed to broadcast cancellation to multiple goroutines
  - `done := make(chan struct{})`
  
- **Pipeline**: Chain of stages connected by channels
  - Each stage is a goroutine reading from inbound and writing to outbound channels
  - First stage = producer, last stage = consumer, middle = transformers
  - Pattern: generator → transformer → consumer

## Key Takeaways
- Goroutines are cheap, spawn many without worry
- Channels are the primary way to communicate between goroutines
- Always close channels when done producing
- Use buffered channels to reduce blocking
- Select enables non-blocking and timeout operations
- Done channels gracefully shut down goroutines
- Pipelines organize concurrent workflows
