package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
BACKPRESSURE
============

What is it?
-----------
Backpressure is a mechanism to handle situations where a producer is faster
than a consumer. It prevents resource exhaustion by controlling the flow of
data through the system.

The Problem:
- Fast producer overwhelms slow consumer
- Unbounded memory growth
- System crashes or degradation
- Message/event loss

Solutions:
1. Buffered Channels: Limited buffer size
2. Blocking: Producer waits for consumer
3. Dropping: Discard excess items (oldest/newest)
4. Rate Limiting: Slow down producer
5. Load Shedding: Reject new work when overloaded
6. Dynamic Batching: Process in batches

Benefits:
- Prevents out-of-memory errors
- System stability under load
- Graceful degradation
- Better resource utilization

Use Cases:
- Message queue systems
- Event streaming
- API request handling
- Data pipelines
- Log aggregation
- Real-time data processing

Best Practices:
1. Set appropriate buffer sizes
2. Monitor queue depth
3. Provide feedback to producer
4. Choose right backpressure strategy
5. Handle dropped items gracefully
6. Use metrics and alerting
7. Test under load
*/

// ============================================================================
// STRATEGY 1: BUFFERED CHANNEL (Basic Backpressure)
// ============================================================================

func producer(ch chan<- int, count int) {
	for i := 1; i <= count; i++ {
		select {
		case ch <- i:
			fmt.Printf("  Producer: Sent %d\n", i)
		default:
			fmt.Printf("  Producer: Channel full! Dropping %d\n", i)
		}
		time.Sleep(50 * time.Millisecond)
	}
	close(ch)
}

func consumer(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for item := range ch {
		fmt.Printf("    Consumer %d: Processing %d...\n", id, item)
		time.Sleep(200 * time.Millisecond) // Slow consumer
		fmt.Printf("    Consumer %d: Finished %d\n", id, item)
	}
}

func demoBufferedChannel() {
	fmt.Println("\n📌 Backpressure: Buffered Channel")
	fmt.Println("-" * 40)

	bufferSize := 3
	ch := make(chan int, bufferSize)

	fmt.Printf("Buffer size: %d\n", bufferSize)
	fmt.Println("Producer: Fast (50ms/item)")
	fmt.Println("Consumer: Slow (200ms/item)\n")

	var wg sync.WaitGroup

	// Start consumer
	wg.Add(1)
	go consumer(1, ch, &wg)

	// Start producer
	go producer(ch, 10)

	wg.Wait()

	fmt.Println("\n✓ Items beyond buffer were dropped (backpressure)")
}

// ============================================================================
// STRATEGY 2: BLOCKING PRODUCER
// ============================================================================

func blockingProducer(ctx context.Context, ch chan<- int, count int) {
	for i := 1; i <= count; i++ {
		select {
		case ch <- i:
			fmt.Printf("  Producer: Sent %d (waited if needed)\n", i)
		case <-ctx.Done():
			fmt.Println("  Producer: Cancelled")
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	close(ch)
}

func demoBlockingProducer() {
	fmt.Println("\n📌 Backpressure: Blocking Producer")
	fmt.Println("-" * 40)

	ch := make(chan int, 2) // Small buffer
	ctx := context.Background()

	fmt.Println("Strategy: Producer blocks when buffer is full")
	fmt.Println("This creates natural backpressure\n")

	var wg sync.WaitGroup

	// Start consumer
	wg.Add(1)
	go consumer(1, ch, &wg)

	// Start blocking producer
	go blockingProducer(ctx, ch, 8)

	wg.Wait()

	fmt.Println("\n✓ Producer slowed down to match consumer (blocking)")
}

// ============================================================================
// STRATEGY 3: DROP OLDEST (Ring Buffer)
// ============================================================================

type RingBuffer struct {
	buffer   []int
	capacity int
	mu       sync.Mutex
	notEmpty *sync.Cond
	closed   bool
}

func NewRingBuffer(capacity int) *RingBuffer {
	rb := &RingBuffer{
		buffer:   make([]int, 0, capacity),
		capacity: capacity,
	}
	rb.notEmpty = sync.NewCond(&rb.mu)
	return rb
}

func (rb *RingBuffer) Push(item int) bool {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if len(rb.buffer) >= rb.capacity {
		// Drop oldest
		dropped := rb.buffer[0]
		rb.buffer = rb.buffer[1:]
		fmt.Printf("    ⚠ Buffer full! Dropped oldest: %d\n", dropped)
	}

	rb.buffer = append(rb.buffer, item)
	rb.notEmpty.Signal()
	return true
}

func (rb *RingBuffer) Pop() (int, bool) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	for len(rb.buffer) == 0 && !rb.closed {
		rb.notEmpty.Wait()
	}

	if len(rb.buffer) == 0 {
		return 0, false
	}

	item := rb.buffer[0]
	rb.buffer = rb.buffer[1:]
	return item, true
}

func (rb *RingBuffer) Close() {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.closed = true
	rb.notEmpty.Broadcast()
}

func demoDropOldest() {
	fmt.Println("\n📌 Backpressure: Drop Oldest (Ring Buffer)")
	fmt.Println("-" * 40)

	rb := NewRingBuffer(3)

	fmt.Println("Strategy: When buffer full, drop oldest items")
	fmt.Println("Buffer capacity: 3\n")

	// Fast producer
	go func() {
		for i := 1; i <= 10; i++ {
			rb.Push(i)
			fmt.Printf("  Producer: Added %d\n", i)
			time.Sleep(50 * time.Millisecond)
		}
		rb.Close()
	}()

	// Slow consumer
	time.Sleep(500 * time.Millisecond) // Delay consumer start

	fmt.Println("\nConsumer starting (delayed)...\n")

	for {
		item, ok := rb.Pop()
		if !ok {
			break
		}
		fmt.Printf("    Consumer: Processing %d\n", item)
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("\n✓ Oldest items dropped to maintain buffer limit")
}

// ============================================================================
// STRATEGY 4: LOAD SHEDDING
// ============================================================================

type LoadShedder struct {
	maxQueueSize int
	queue        chan int
	processed    int
	dropped      int
	mu           sync.Mutex
}

func NewLoadShedder(maxQueueSize int) *LoadShedder {
	return &LoadShedder{
		maxQueueSize: maxQueueSize,
		queue:        make(chan int, maxQueueSize),
	}
}

func (ls *LoadShedder) Submit(item int) bool {
	select {
	case ls.queue <- item:
		return true
	default:
		ls.mu.Lock()
		ls.dropped++
		ls.mu.Unlock()
		return false
	}
}

func (ls *LoadShedder) Process(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for item := range ls.queue {
		fmt.Printf("    Worker %d: Processing %d\n", id, item)
		time.Sleep(150 * time.Millisecond)

		ls.mu.Lock()
		ls.processed++
		ls.mu.Unlock()
	}
}

func (ls *LoadShedder) Stats() (processed, dropped int) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	return ls.processed, ls.dropped
}

func demoLoadShedding() {
	fmt.Println("\n📌 Backpressure: Load Shedding")
	fmt.Println("-" * 40)

	shedder := NewLoadShedder(5)

	fmt.Println("Strategy: Reject new work when queue is full")
	fmt.Println("Max queue size: 5\n")

	var wg sync.WaitGroup

	// Start workers
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go shedder.Process(i, &wg)
	}

	// Submit work
	submitted := 0
	for i := 1; i <= 20; i++ {
		if shedder.Submit(i) {
			submitted++
			fmt.Printf("  ✓ Submitted %d\n", i)
		} else {
			fmt.Printf("  ✗ Rejected %d (load shedding)\n", i)
		}
		time.Sleep(30 * time.Millisecond)
	}

	close(shedder.queue)
	wg.Wait()

	processed, dropped := shedder.Stats()
	fmt.Printf("\nStats: Submitted=%d, Processed=%d, Dropped=%d\n",
		submitted, processed, dropped)
	fmt.Println("✓ Load shedding prevented queue overflow")
}

// ============================================================================
// STRATEGY 5: DYNAMIC BATCHING
// ============================================================================

func dynamicBatcher(input <-chan int, output chan<- []int, batchSize int, timeout time.Duration) {
	batch := make([]int, 0, batchSize)
	timer := time.NewTimer(timeout)
	timer.Stop()

	flush := func() {
		if len(batch) > 0 {
			output <- batch
			batch = make([]int, 0, batchSize)
		}
		timer.Stop()
	}

	for {
		select {
		case item, ok := <-input:
			if !ok {
				flush()
				close(output)
				return
			}

			batch = append(batch, item)

			if len(batch) == 1 {
				timer.Reset(timeout)
			}

			if len(batch) >= batchSize {
				flush()
			}

		case <-timer.C:
			flush()
		}
	}
}

func demoDynamicBatching() {
	fmt.Println("\n📌 Backpressure: Dynamic Batching")
	fmt.Println("-" * 40)

	input := make(chan int, 10)
	output := make(chan []int)

	fmt.Println("Strategy: Batch items to reduce processing overhead")
	fmt.Println("Batch size: 5, Timeout: 200ms\n")

	// Start batcher
	go dynamicBatcher(input, output, 5, 200*time.Millisecond)

	// Consumer
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for batch := range output {
			fmt.Printf("  Processing batch of %d items: %v\n", len(batch), batch)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Producer
	go func() {
		for i := 1; i <= 13; i++ {
			input <- i
			fmt.Printf("    Produced %d\n", i)
			time.Sleep(50 * time.Millisecond)
		}
		close(input)
	}()

	wg.Wait()

	fmt.Println("\n✓ Batching reduces per-item overhead and backpressure")
}

// ============================================================================
// PRACTICAL EXAMPLE: EVENT STREAM PROCESSOR
// ============================================================================

type Event struct {
	ID        int
	Timestamp time.Time
	Data      string
}

type EventProcessor struct {
	inbound   chan Event
	processed chan Event
	dropped   int
	mu        sync.Mutex
}

func NewEventProcessor(bufferSize int) *EventProcessor {
	return &EventProcessor{
		inbound:   make(chan Event, bufferSize),
		processed: make(chan Event, bufferSize),
	}
}

func (ep *EventProcessor) Receive(event Event) bool {
	select {
	case ep.inbound <- event:
		return true
	default:
		ep.mu.Lock()
		ep.dropped++
		ep.mu.Unlock()
		return false
	}
}

func (ep *EventProcessor) Start(workers int) {
	var wg sync.WaitGroup

	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for event := range ep.inbound {
				// Simulate processing
				time.Sleep(100 * time.Millisecond)
				ep.processed <- event
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(ep.processed)
	}()
}

func (ep *EventProcessor) DroppedCount() int {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	return ep.dropped
}

func demoPracticalExample() {
	fmt.Println("\n📌 Practical Example: Event Stream Processor")
	fmt.Println("-" * 40)

	processor := NewEventProcessor(5)
	processor.Start(2)

	fmt.Println("Event processor with backpressure control")
	fmt.Println("Buffer: 5 events, Workers: 2\n")

	// Event generator (fast)
	go func() {
		for i := 1; i <= 15; i++ {
			event := Event{
				ID:        i,
				Timestamp: time.Now(),
				Data:      fmt.Sprintf("Event-%d", i),
			}

			if processor.Receive(event) {
				fmt.Printf("  ✓ Event %d received\n", i)
			} else {
				fmt.Printf("  ✗ Event %d dropped (backpressure)\n", i)
			}

			time.Sleep(40 * time.Millisecond)
		}
		close(processor.inbound)
	}()

	// Process events
	processed := 0
	for event := range processor.processed {
		processed++
		fmt.Printf("    Processed: %s\n", event.Data)
	}

	dropped := processor.DroppedCount()
	fmt.Printf("\nFinal Stats: Processed=%d, Dropped=%d\n", processed, dropped)
	fmt.Println("✓ Backpressure prevented system overload")
}

// ============================================================================
// BEST PRACTICES
// ============================================================================

func demoBestPractices() {
	fmt.Println("\n📌 Backpressure Best Practices")
	fmt.Println("-" * 40)

	fmt.Println("\n✅ Choose Right Strategy:")
	fmt.Println("  • Blocking: When all data is critical")
	fmt.Println("  • Drop Oldest: Time-sensitive data (metrics)")
	fmt.Println("  • Drop Newest: Priority to established data")
	fmt.Println("  • Load Shedding: Protect system under load")
	fmt.Println("  • Batching: Reduce overhead")

	fmt.Println("\n✅ Monitoring:")
	fmt.Println("  • Track queue depth")
	fmt.Println("  • Monitor drop rates")
	fmt.Println("  • Alert on sustained pressure")
	fmt.Println("  • Log backpressure events")

	fmt.Println("\n✅ Configuration:")
	fmt.Println("  • Make buffer sizes configurable")
	fmt.Println("  • Set appropriate timeouts")
	fmt.Println("  • Tune based on load testing")
	fmt.Println("  • Consider auto-scaling")

	fmt.Println("\n✅ Error Handling:")
	fmt.Println("  • Provide feedback on drops")
	fmt.Println("  • Log dropped items (if critical)")
	fmt.Println("  • Consider retry mechanisms")
	fmt.Println("  • Graceful degradation")
}

// ============================================================================
// DEMONSTRATION
// ============================================================================

// DemoBackpressure demonstrates backpressure handling patterns
func DemoBackpressure() {
	fmt.Println("\n" + "="*60)
	fmt.Println("BACKPRESSURE DEMONSTRATION")
	fmt.Println("="*60 + "\n")

	demoBufferedChannel()
	demoBlockingProducer()
	demoDropOldest()
	demoLoadShedding()
	demoDynamicBatching()
	demoPracticalExample()
	demoBestPractices()

	fmt.Println("\n✅ Backpressure demonstration completed!")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("- Backpressure prevents resource exhaustion")
	fmt.Println("- Choose strategy based on requirements")
	fmt.Println("- Buffered channels provide natural backpressure")
	fmt.Println("- Monitor queue depth and drop rates")
	fmt.Println("- Load shedding protects system under stress")
	fmt.Println("- Batching can reduce backpressure")
}
