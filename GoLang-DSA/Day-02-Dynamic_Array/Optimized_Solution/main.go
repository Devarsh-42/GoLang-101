package main

import (
	"fmt"
	"math/rand"
	"time"
)

// a video data packet
type Packet struct {
	SequenceID int
	Data       int
	Timestamp  int64
}

// ring buffer for video streaming — no compaction needed
type StreamBuffer struct {
	packets  []Packet
	head     int 
	tail     int 
	count    int 
	capacity int 
	dropped  int 
}

func NewStreamBuffer(capacity int) *StreamBuffer {
	if capacity <= 0 {
		capacity = 1
	}
	return &StreamBuffer{
		packets:  make([]Packet, capacity),
		head:     0,
		tail:     0,
		count:    0,
		capacity: capacity,
		dropped:  0,
	}
}

func (sb *StreamBuffer) Push(p Packet) {
	if sb.count == sb.capacity {
		fmt.Printf("  [OVERFLOW] Buffer full! Dropping oldest packet #%d to make room\n",
			sb.packets[sb.head].SequenceID)
		sb.head = (sb.head + 1) % sb.capacity // advance head, discard oldest
		sb.count--
		sb.dropped++
	}
	sb.packets[sb.tail] = p
	sb.tail = (sb.tail + 1) % sb.capacity // wrap around
	sb.count++
}

func (sb *StreamBuffer) Play() (Packet, bool) {
	if sb.count == 0 {
		return Packet{}, false // buffer underrun! (buffering...)
	}
	p := sb.packets[sb.head]
	sb.packets[sb.head] = Packet{}        // zero out for cleanliness
	sb.head = (sb.head + 1) % sb.capacity // wrap around
	sb.count--
	return p, true
}

func (sb *StreamBuffer) UnreadCount() int {
	return sb.count
}

func (sb *StreamBuffer) Stats() {
	fmt.Println("\n========= RING BUFFER STATS =========")
	fmt.Printf("Capacity (fixed):  %d\n", sb.capacity)
	fmt.Printf("Buffered (unread): %d\n", sb.count)
	fmt.Printf("Head (read pos):   %d\n", sb.head)
	fmt.Printf("Tail (write pos):  %d\n", sb.tail)
	fmt.Printf("Dropped packets:   %d\n", sb.dropped)
	if sb.capacity > 0 {
		fmt.Printf("Utilization:       %d%%\n", (sb.count*100)/sb.capacity)
	}
	fmt.Println("=====================================")
}

func main() {
	buffer := NewStreamBuffer(8) 

	totalPackets := rand.Intn(21) + 20 // 20-40 packets
	fmt.Printf("Video Stream Simulation (Ring Buffer) | %d packets incoming | fixed capacity: 8\n\n", totalPackets)

	playedCount := 0

	for i := 0; i < totalPackets; i++ { // Simulate packet arrival
		p := Packet{
			SequenceID: i + 1,
			Data:       rand.Intn(256),
			Timestamp:  time.Now().UnixMicro(),
		}
		fmt.Printf("Recv packet #%d (data=%d)\n", p.SequenceID, p.Data)
		buffer.Push(p)

		// Every 3 packets received -> try to play 2 (simulates playback)
		if (i+1)%3 == 0 {
			for j := 0; j < 2; j++ {
				pkt, ok := buffer.Play()
				if ok {
					fmt.Printf("  Playing packet #%d (data=%d) | buffered ahead: %d\n",
						pkt.SequenceID, pkt.Data, buffer.UnreadCount())
					playedCount++
				} else {
					fmt.Println("  BUFFERING... no packets available!")
				}
			}
		}
	}

	// Drain remaining buffer
	fmt.Println("\n--- Draining remaining buffer ---")
	for {
		pkt, ok := buffer.Play()
		if !ok {
			break
		}
		fmt.Printf("  Playing packet #%d (data=%d)\n", pkt.SequenceID, pkt.Data)
		playedCount++
	}

	fmt.Printf("\nTotal packets played: %d\n", playedCount)
	buffer.Stats()
}
