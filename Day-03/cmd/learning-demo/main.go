package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"day03/learning"
)

// ==========================================
// DAY 03 - LEARNING DEMO
// ==========================================
// This program demonstrates:
// 1. Maps (key-value data structures)
// 2. Structs and Custom Types
// 3. Receiver Functions (Methods)

func main() {
	fmt.Println("╔════════════════════════════════════════════════════╗")
	fmt.Println("║                  DAY 03 - GO LEARNING              ║")
	fmt.Println("║   Maps, Structs, Custom Types & Receiver Functions ║")
	fmt.Println("╚════════════════════════════════════════════════════╝")

	// Interactive menu
	reader := bufio.NewReader(os.Stdin)

	for {
		displayMenu()
		fmt.Print("\nEnter your choice (1-4, or 'q' to quit): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			learning.RunAllMapDemos()
			pressEnterToContinue(reader)

		case "2":
			learning.RunAllStructDemos()
			pressEnterToContinue(reader)

		case "3":
			learning.RunAllReceiverFunctionDemos()
			pressEnterToContinue(reader)

		case "4":
			runAllDemos()
			pressEnterToContinue(reader)

		case "q", "Q":
			fmt.Println("\n👋 Thanks for learning Go! Happy coding!")
			return

		default:
			fmt.Println("\n❌ Invalid choice. Please enter 1-4 or 'q'")
		}
	}
}

func displayMenu() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("              LEARNING MENU")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("1. Learn about MAPS")
	fmt.Println("   - Map basics (create, read, update, delete)")
	fmt.Println("   - Map iteration")
	fmt.Println("   - Complex maps (nested, with structs)")
	fmt.Println("   - Practical examples")
	fmt.Println()
	fmt.Println("2. Learn about STRUCTS & CUSTOM TYPES")
	fmt.Println("   - Basic structs")
	fmt.Println("   - Nested structs")
	fmt.Println("   - Anonymous structs")
	fmt.Println("   - Custom types")
	fmt.Println("   - Struct tags and pointers")
	fmt.Println()
	fmt.Println("3. Learn about RECEIVER FUNCTIONS (Methods)")
	fmt.Println("   - Value receivers vs Pointer receivers")
	fmt.Println("   - Methods on custom types")
	fmt.Println("   - Methods with parameters")
	fmt.Println("   - Method chaining")
	fmt.Println()
	fmt.Println("4. Run ALL Demonstrations")
	fmt.Println()
	fmt.Println("Q. Quit")
	fmt.Println(strings.Repeat("=", 50))
}

func runAllDemos() {
	fmt.Println("\n" + strings.Repeat("*", 50))
	fmt.Println("      RUNNING ALL DEMONSTRATIONS")
	fmt.Println(strings.Repeat("*", 50))

	learning.RunAllMapDemos()
	fmt.Println("\n" + strings.Repeat("-", 50))

	learning.RunAllStructDemos()
	fmt.Println("\n" + strings.Repeat("-", 50))

	learning.RunAllReceiverFunctionDemos()
	fmt.Println("\n" + strings.Repeat("*", 50))
}

func pressEnterToContinue(reader *bufio.Reader) {
	fmt.Print("\n\nPress Enter to return to menu...")
	reader.ReadString('\n')
	clearScreen()
}

func clearScreen() {
	// Simple way to add some space
	fmt.Println("\n\n")
}

// ==========================================
// QUICK REFERENCE GUIDE
// ==========================================
/*

MAPS:
-----
// Create
m := make(map[string]int)
m := map[string]int{"key": value}

// Operations
m[key] = value           // Add/Update
value := m[key]          // Get
value, ok := m[key]      // Get with existence check
delete(m, key)           // Delete
len(m)                   // Size

// Iterate
for key, value := range m { }


STRUCTS:
--------
// Define
type Person struct {
    Name string
    Age  int
}

// Create
p := Person{Name: "John", Age: 30}
p := Person{"John", 30}
p := &Person{Name: "John"}  // Pointer

// Access
p.Name = "Jane"
age := p.Age


CUSTOM TYPES:
-------------
type UserID int
type Status string

var id UserID = 100


RECEIVER FUNCTIONS (METHODS):
------------------------------
// Value receiver (copy)
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Pointer receiver (original)
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

// Call
rect := Rectangle{Width: 10, Height: 5}
area := rect.Area()
rect.Scale(2)

*/
