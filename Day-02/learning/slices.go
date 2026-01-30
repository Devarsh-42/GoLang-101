package learning

import (
	"fmt"
	"slices"
)

// Function to demonstrate slice usage
func ExampleSlices() {
	// Creating a slice
	numbers := []int{10, 20, 30, 40, 50}
	fmt.Printf("Original slice: %v\n", numbers)

	// Appending to a slice
	numbers = append(numbers, 60, 70)
	fmt.Printf("After appending: %v\n", numbers)

	// Slicing a slice
	subSlice := numbers[2:5]
	fmt.Printf("Sliced portion (index 2 to 4): %v\n", subSlice)

	// Using slices package to sort
	slices.Sort(numbers)
	fmt.Printf("Sorted slice: %v\n", numbers)

	// Checking if a slice contains an element
	contains := slices.Contains(numbers, 30)
	fmt.Printf("Slice contains 30: %v\n", contains)


	// Inefficient approach
	var result []int
	for i := 0; i < 1000; i++ {
		result = append(result, i)  // Multiple reallocations
	}

	// Efficient approach
	result2 := make([]int, 0, 1000)  // Preallocate capacity
	for i := 0; i < 1000; i++ {
		result2 = append(result2, i)  // Minimal reallocation
	}

	// Error Handling - Preventing Slice Bounds Errors

	// 	func safeAccess(slice []int, index int) (int, error) {
	// 	if index < 0 || index >= len(slice) {
	// 		return 0, fmt.Errorf("index out of bounds")
	// 	}
	// 	return slice[index], nil
	// }

	// Passing slice to a function
	modifySlice(numbers)
	fmt.Printf("After modifying slice in function: %v\n", numbers)
}

// passing slices to functions
func modifySlice(s []int) {
	for i := range s {
		s[i] *= 2
	}
}