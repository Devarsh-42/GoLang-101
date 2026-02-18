package main

import (
	"day02/learning"
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("=== Go Learning Examples ===\n")

	// Structs examples
	fmt.Println("1. Structs:")
	learning.ExampleStructs()
	fmt.Println()

	// Struct methods
	fmt.Println("2. Struct Methods:")
	learning.ExampleStructMethods()
	fmt.Println()

	// Arrays examples
	fmt.Println("3. Arrays:")
	learning.ExampleArrays()

	// Error Examples
	fmt.Println("\n4. Errors:")
	learning.ExampleErrors()

	// Interface Examples
	fmt.Println("\n6. Interfaces:")
	learning.ExampleInterfaces()

	//Slices Examples
	fmt.Println("\n5. Slices:")
	learning.ExampleSlices()

	// Loops Examples
	fmt.Println("\n7. Loops:")
	learning.ExampleLoops()
	runtime.GC()

}
