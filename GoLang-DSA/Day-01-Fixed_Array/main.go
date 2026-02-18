package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

// Theater holds the booking state
type Theater struct {
	mu    sync.RWMutex // RWMutex: multiple readers OR one writer
	name  string // theater name
	seats []bool
}

// NewTheater is a constructor
// Why? It guarantees a Theater is never in an invalid state.
// Without this, someone could do Theater{} and have 0 seats.
func NewTheater(name string, totalSeats int) *Theater {
	if totalSeats <= 0 {
		panic("theater must have at least 1 seat")
	}
	return &Theater{
		name:  name,
		seats: make([]bool, totalSeats),
	}
}

func (t *Theater) TotalSeats() int {
	// No lock needed — len(t.seats) never changes after construction
	return len(t.seats)
}

// ─── READ operations use RLock (multiple goroutines can read simultaneously) ───

func (t *Theater) displayAvailableSeats() {
	t.mu.RLock() // multiple readers allowed at the same time
	defer t.mu.RUnlock()

	fmt.Println("\n--- Available Seats ---")
	count := 0
	for i, booked := range t.seats {
		if !booked {
			fmt.Printf("%d ", i+1)
			count++
		}
	}
	fmt.Printf("\n\nTotal available: %d / %d\n", count, t.TotalSeats())
}

func (t *Theater) isSeatAvailable(seatNo int) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return !t.seats[seatNo-1]
}

func (t *Theater) countAvailable() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	count := 0
	for _, booked := range t.seats {
		if !booked {
			count++
		}
	}
	return count
}

// ─── WRITE operations use Lock (exclusive — blocks all readers and writers) ───

func (t *Theater) bookSeat() {
	seatNo, ok := t.readSeatNumber()
	if !ok {
		return
	}

	t.mu.Lock() // Exclusive lock — no one else reads or writes
	defer t.mu.Unlock()

	index := seatNo - 1
	if t.seats[index] {
		fmt.Printf("Seat %d is already booked.\n", seatNo)
		return
	}

	t.seats[index] = true
	fmt.Printf("Seat %d booked successfully!\n", seatNo)
}

func (t *Theater) cancelBooking() {
	seatNo, ok := t.readSeatNumber()
	if !ok {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	index := seatNo - 1
	if !t.seats[index] {
		fmt.Printf("Seat %d is not booked.\n", seatNo)
		return
	}

	t.seats[index] = false
	fmt.Printf("Seat %d booking cancelled.\n", seatNo)
}

// bookSeatConcurrent is used by the simulation (no user input)
func (t *Theater) bookSeatConcurrent(seatNo int, userID int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Simulate some processing time while holding the lock
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)

	index := seatNo - 1
	if t.seats[index] {
		fmt.Printf("REJECTED - [User %d] Seat %d is already booked.\n", userID, seatNo)
		return
	}

	t.seats[index] = true
	fmt.Printf("  User %d] Seat %d booked successfully!\n", userID, seatNo)
}

// checkSeatConcurrent is used by the simulation (read-only)
func (t *Theater) checkSeatConcurrent(seatNo int, userID int) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Simulate some read processing time
	time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)

	status := "available"
	if t.seats[seatNo-1] {
		status = "booked"
	}
	fmt.Printf(" [User %d] Checked seat %d → %s\n", userID, seatNo, status)
}

func (t *Theater) readSeatNumber() (int, bool) {
	var seatNo int
	fmt.Printf("Enter seat number (1-%d): ", t.TotalSeats())
	if _, err := fmt.Scan(&seatNo); err != nil {
		fmt.Println("Invalid input.")
		return 0, false
	}

	if seatNo < 1 || seatNo > t.TotalSeats() {
		fmt.Printf("Invalid seat number. Must be between 1 and %d.\n", t.TotalSeats())
		return 0, false
	}

	return seatNo, true
}

// ─── Concurrent Simulation ───

func simulateConcurrency(theater *Theater) {
	fmt.Println("\n═══ CONCURRENT BOOKING SIMULATION ═══")
	fmt.Printf("   Theater: %s | Seats: %d\n", theater.name, theater.TotalSeats())
	fmt.Println("   Spawning 20 users trying to book 10 seats...")
	fmt.Println("   (This means collisions WILL happen)")

	var wg sync.WaitGroup

	// Phase 1: Multiple readers checking seats simultaneously
	// RLock allows ALL of these to run in parallel
	fmt.Println("── Phase 1: Concurrent READS (all run in parallel with RLock) ──")
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			seatNo := rand.Intn(theater.TotalSeats()) + 1
			theater.checkSeatConcurrent(seatNo, userID)
		}(i)
	}
	wg.Wait()

	// Phase 2: Multiple writers trying to book the SAME seats
	// Lock ensures only ONE writer at a time — prevents double booking
	fmt.Println("\n── Phase 2: Concurrent WRITES (serialized with Lock) ──")
	fmt.Println("   20 users fighting over 10 seats...")
	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			seatNo := rand.Intn(10) + 1 // only seats 1-10, forcing collisions
			theater.bookSeatConcurrent(seatNo, userID)
		}(i)
	}
	wg.Wait()

	// Phase 3: Mixed reads and writes simultaneously
	// Readers use RLock (parallel), writers use Lock (exclusive)
	// A writer BLOCKS all readers until done
	fmt.Println("\n── Phase 3: Mixed READS + WRITES (RLock vs Lock contention) ──")
	for i := 1; i <= 15; i++ {
		wg.Add(1)
		if i%3 == 0 {
			// Every 3rd user is a writer
			go func(userID int) {
				defer wg.Done()
				seatNo := rand.Intn(theater.TotalSeats()) + 1
				theater.bookSeatConcurrent(seatNo, userID)
			}(100 + i) // userID 100+ to distinguish
		} else {
			// Others are readers
			go func(userID int) {
				defer wg.Done()
				seatNo := rand.Intn(theater.TotalSeats()) + 1
				theater.checkSeatConcurrent(seatNo, userID)
			}(200 + i)
		}
	}
	wg.Wait()

	// Final summary
	fmt.Println("\n── Final State ──")
	available := theater.countAvailable()
	booked := theater.TotalSeats() - available
	fmt.Printf("   Booked: %d | Available: %d | Total: %d\n", booked, available, theater.TotalSeats())
	fmt.Println("\n  No double bookings. No race conditions. RWMutex worked.")
	fmt.Println("═══════════════════════════════════════════")
}

func main() {
	var numSeats int
	fmt.Print("Enter total number of seats: ")
	if _, err := fmt.Scan(&numSeats); err != nil || numSeats <= 0 {
		fmt.Println("Invalid number of seats.")
		os.Exit(1)
	}

	theater := NewTheater("Main Hall", numSeats)

	for {
		fmt.Println("\n------- Seat Booking Program -------")
		fmt.Printf("Theater: %s (%d seats)\n", theater.name, theater.TotalSeats())
		fmt.Println("1. View available seats")
		fmt.Println("2. Book a seat")
		fmt.Println("3. Cancel a booking")
		fmt.Println("4. Run concurrent simulation")
		fmt.Println("5. Exit")
		fmt.Print("Enter your choice: ")

		var choice int
		if _, err := fmt.Scan(&choice); err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch choice {
		case 1:
			theater.displayAvailableSeats()
		case 2:
			theater.bookSeat()
		case 3:
			theater.cancelBooking()
		case 4:
			simulateConcurrency(theater)
		case 5:
			fmt.Println("Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please select 1-5.")
		}
	}
}
