package main

import (
	"fmt"
	"math/rand"
)

type VideoBuffer struct {
	packets  []int 
	size     int   
	capacity int   
}

func NewVideoBuffer(initialCapacity int) *VideoBuffer {
	if initialCapacity <= 0 {
		initialCapacity = 1
	}
	return &VideoBuffer{
		packets:  make([]int, initialCapacity),
		size:     0,
		capacity: initialCapacity,
	}
}

func (b *VideoBuffer) Add(packet int) {
	if b.size == b.capacity {
		newCapacity := b.capacity * 2
		newPackets := make([]int, newCapacity)

		for i := 0; i < b.size; i++ {
			newPackets[i] = b.packets[i]
		}

		b.packets = newPackets
		b.capacity = newCapacity
		fmt.Printf("  >> Buffer FULL! Doubled capacity: %d -> %d\n", b.capacity/2, b.capacity)
	}

	b.packets[b.size] = packet
	b.size++
}

func (b *VideoBuffer) Display() {
	fmt.Printf("Buffer [size=%d, capacity=%d]: ", b.size, b.capacity)
	for i := 0; i < b.size; i++ {
		fmt.Printf("%d ", b.packets[i])
	}
	fmt.Println()
}

func main() {
	totalPackets := rand.Intn(21) + 10 // 10 to 30 packets

	fmt.Printf("Simulating %d incoming video packets...\n", totalPackets)
	fmt.Println("Starting buffer with initial capacity: 4")
	fmt.Println()

	buffer := NewVideoBuffer(4)

	for i := 0; i < totalPackets; i++ {
		packet := rand.Intn(256) // simulate packet data (0-255)
		fmt.Printf("Packet #%d arrived: %d\n", i+1, packet)
		buffer.Add(packet)
	}

	fmt.Println()
	buffer.Display()
}
