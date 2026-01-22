package main

import (
	"fmt"

	"day01/learning"
)

func main() {
	fmt.Println("=== Go Learning Examples ===\n")

	// Variables examples
	fmt.Println("1. Variables:")
	x := 42                // short declaration
	var name string = "Go" // explicit type
	const maxValue = 100   // constant
	fmt.Printf("   x = %d, name = %s, maxValue = %d\n\n", x, name, maxValue)

	// Control flow - if/else
	fmt.Println("2. Control Flow - If/Else:")
	score := 85
	if score >= 90 {
		fmt.Println("   Grade: A")
	} else if score >= 80 {
		fmt.Println("   Grade: B")
	} else {
		fmt.Println("   Grade: C")
	}
	fmt.Println()

	learning.ExampleIfElse()
	fmt.Println()

	// For loop
	fmt.Println("3. For Loop:")
	for i := 1; i <= 5; i++ {
		fmt.Printf("   Count: %d\n", i)
	}
	fmt.Println()

	// Range over slice
	fmt.Println("4. Range over Slice:")
	numbers := []int{10, 20, 30, 40, 50}
	for i, v := range numbers {
		fmt.Printf("   Index %d: Value %d\n", i, v)
	}
	fmt.Println()

	// Switch statement
	fmt.Println("5. Switch Statement:")
	day := "Monday"
	switch day {
	case "Monday":
		fmt.Println("   Start of the work week")
	case "Friday":
		fmt.Println("   Almost weekend!")
	default:
		fmt.Println("   Mid-week")
	}
	fmt.Println()

	// Functions
	fmt.Println("6. Functions:")
	result := add(5, 3)
	fmt.Printf("   add(5, 3) = %d\n", result)

	q, r := divideWithRemainder(17, 5)
	fmt.Printf("   17 / 5 = %d remainder %d\n", q, r)
	fmt.Println()

	// Defer
	fmt.Println("7. Defer (executes in LIFO order):")
	deferExample()
	fmt.Println()

	fmt.Println("=== Examples Complete ===")
}

// Function with single return
func add(a, b int) int {
	return a + b
}

// Function with multiple returns
func divideWithRemainder(a, b int) (quotient, remainder int) {
	quotient = a / b
	remainder = a % b
	return
}

// Defer example
func deferExample() {
	defer fmt.Println("   Third (deferred)")
	defer fmt.Println("   Second (deferred)")
	defer fmt.Println("   First (deferred)")
	fmt.Println("   Regular execution")
}
