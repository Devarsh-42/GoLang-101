package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
GO CONCURRENCY PATTERNS - COMPREHENSIVE DEMONSTRATION
======================================================

This program demonstrates practical implementations of all major
Go concurrency patterns and best practices.

Topics Covered:
✅ Worker Pools
✅ Fan-out / Fan-in
✅ Select
✅ Context (cancellation & timeout)
✅ sync.WaitGroup
✅ sync.Mutex / RWMutex
✅ Atomic operations
✅ Rate limiting
✅ Deadlock prevention
✅ Graceful shutdown
✅ Backpressure

Author: Devarsh Mehta
Date: Day-08 of GoLang-101
*/

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                               ║")
	fmt.Println("║        GO CONCURRENCY PATTERNS DEMONSTRATION                  ║")
	fmt.Println("║        Complete Guide to Go Concurrency                       ║")
	fmt.Println("║                                                               ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")

	fmt.Println("\n📚 This demonstration covers 11 essential concurrency topics")
	fmt.Println("   Each topic includes explanations, examples, and best practices")

	topics := []struct {
		id          int
		name        string
		description string
		demo        func()
	}{
		{
			id:          1,
			name:        "Worker Pools",
			description: "Fixed goroutines processing jobs from a queue",
			demo:        DemoWorkerPools,
		},
		{
			id:          2,
			name:        "Fan-out / Fan-in",
			description: "Distribute work and merge results",
			demo:        DemoFanOutFanIn,
		},
		{
			id:          3,
			name:        "Select Statement",
			description: "Multiplexing channel operations",
			demo:        DemoSelect,
		},
		{
			id:          4,
			name:        "Context Package",
			description: "Cancellation and timeout management",
			demo:        DemoContext,
		},
		{
			id:          5,
			name:        "sync.WaitGroup",
			description: "Wait for goroutines to complete",
			demo:        DemoWaitGroup,
		},
		{
			id:          6,
			name:        "Mutex & RWMutex",
			description: "Protecting shared state",
			demo:        DemoMutex,
		},
		{
			id:          7,
			name:        "Atomic Operations",
			description: "Lock-free synchronization",
			demo:        DemoAtomic,
		},
		{
			id:          8,
			name:        "Rate Limiting",
			description: "Controlling operation rate",
			demo:        DemoRateLimiting,
		},
		{
			id:          9,
			name:        "Deadlock Prevention",
			description: "Avoiding circular dependencies",
			demo:        DemoDeadlockPrevention,
		},
		{
			id:          10,
			name:        "Graceful Shutdown",
			description: "Clean application termination",
			demo:        DemoGracefulShutdown,
		},
		{
			id:          11,
			name:        "Backpressure",
			description: "Handling fast producers",
			demo:        DemoBackpressure,
		},
	}

	// Interactive menu
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n" + strings.Repeat("━", 65))
		fmt.Println("📋 MENU - Select a topic to demonstrate:")
		fmt.Println(strings.Repeat("━", 65))

		for _, topic := range topics {
			fmt.Printf("  %2d. %-25s - %s\n", topic.id, topic.name, topic.description)
		}

		fmt.Println("\n   0. Run All Demonstrations")
		fmt.Println("   Q. Quit")

		fmt.Print("\n👉 Enter your choice: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Handle quit
		if strings.ToLower(input) == "q" || strings.ToLower(input) == "quit" {
			fmt.Println("\n✅ Thank you for exploring Go concurrency patterns!")
			fmt.Println("🚀 Happy coding!")
			break
		}

		// Handle run all
		if input == "0" {
			fmt.Println("\n🚀 Running all demonstrations...")
			fmt.Println("⚠️  This will take several minutes. Press Ctrl+C to stop.\n")

			for _, topic := range topics {
				fmt.Println("\n" + strings.Repeat("═", 70))
				fmt.Printf("Running: %s\n", topic.name)
				fmt.Println(strings.Repeat("═", 70))
				topic.demo()

				// Pause between demos
				fmt.Println("\n⏸  Press Enter to continue to next topic...")
				reader.ReadString('\n')
			}

			fmt.Println("\n🎉 All demonstrations completed!")
			continue
		}

		// Parse choice
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(topics) {
			fmt.Println("\n❌ Invalid choice. Please enter a number from the menu.")
			continue
		}

		// Run selected demo
		selected := topics[choice-1]

		fmt.Println("\n" + strings.Repeat("═", 70))
		fmt.Printf("Running: %s\n", selected.name)
		fmt.Printf("Description: %s\n", selected.description)
		fmt.Println(strings.Repeat("═", 70))

		selected.demo()

		fmt.Println("\n✅ Demonstration completed!")
		fmt.Println("📝 Review the code in the corresponding file for detailed comments")
		fmt.Print("\nPress Enter to return to menu...")
		reader.ReadString('\n')
	}
}
