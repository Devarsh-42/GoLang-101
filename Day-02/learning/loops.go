package learning

import (
	"fmt"
)

// Loops in Go
func ExampleLoops() {
	// For loop
	fmt.Println("For Loop:")
	for i := 0; i < 5; i++ {
		fmt.Printf("Iteration %d\n", i)
	}

	// While-like loop
	fmt.Println("\nWhile-like Loop:")
	count := 0
	for count < 5 {
		fmt.Printf("Count is %d\n", count)
		count++
	}

	// // Infinite loop with break
	// fmt.Println("\nInfinite Loop with Break:")
	// i := 0
	// for {
	// 	if i >= 5 {
	// 		break
	// 	}
	// 	fmt.Printf("i is %d\n", i)
	// 	i++
	// }

	// Looping over a slice
	fmt.Println("\nLooping over a Slice:")
	numbers := []int{10, 20, 30, 40, 50}
	for index, value := range numbers {
		fmt.Printf("Index: %d, Value: %d\n", index, value)
	}

	// Looping over a map
	fmt.Println("\nLooping over a Map:")
	person := map[string]string{"name": "Alice", "city": "Wonderland"}
	for key, value := range person {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}
}
