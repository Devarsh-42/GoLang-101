package learning

import "fmt"

// Arrays example
func ExampleArrays() {
	// Declare and initialize an array of integers
	var numbers [5]int = [5]int{10, 20, 30, 40, 50}

	// Declare and initialize an array with inferred size
	cities := [...]string{"New York", "Los Angeles", "Chicago"}

	// Print arrays
	fmt.Println("Numbers array:", numbers)
	fmt.Println("Cities array:", cities)

	// Declare an array without initialization (default values)
	var flags [3]bool

	fmt.Println("Flags array (default values):", flags)

	// Access and modify array elements
	numbers[0] = 15

	// Print array elements
	for i, v := range numbers {
		fmt.Printf("Index %d: Value %d\n", i, v)
	}

	// Multi-dimensional array
	var matrix [2][3]int = [2][3]int{
		{1, 2, 3},
		{4, 5, 6},
	}

	// Print multi-dimensional array elements
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			fmt.Printf("matrix[%d][%d] = %d\n", i, j, matrix[i][j])
		}
	}
}
