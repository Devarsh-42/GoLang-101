package learning

import (
	"fmt"
	"sort"
)

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

	// Array Slicing in Golang
	arr := [3]int{20,30,40}
	fmt.Println(arr[0:2]) // [1 2]
	fmt.Println(arr[:2]) // [1 2]
	fmt.Println(arr[2:]) // [3 4]
	fmt.Println(arr[:]) // [1 2 3 4]

	// Access and modify array elements
	numbers[0] = 15

	// Print array elements
	for i, v := range numbers {
		fmt.Printf("Index %d: Value %d\n", i, v)
	}

	// Array Comparision 
	arr1 := [4]int{1, 2, 3, 4}
	arr2 := [4]int{1, 2, 3, 4}
	arr3 := [4]int{1, 2, 3, 5}
	fmt.Println(arr1 == arr2) // true
	fmt.Println(arr1 == arr3) // false


	// Array Sorting in Golang

	arr4 := [4]int{4, 3, 2, 1}
    sort.Ints(arr4[:]) // ->function that sorts an array of integers in ascending order
    fmt.Println(arr4) // [1 2 3 4]

	// Array Copy in Golang

	// arr1 := [4]int{1, 2, 3, 4}
	// arr2 := [4]int{5, 6, 7, 8}
	// n := copy(arr2[:], arr1[:])
	// fmt.Println(arr2) // [1 2 3 4]
	// fmt.Println(n) // 4


	// Array Contains in Golang

	// arr := [4]int{1, 2, 3, 4}
	// i := sort.Search(len(arr), func(i int) bool { return arr[i] >= 3 })
	// if i < len(arr) && arr[i] == 3 {
	// 	fmt.Println("found 3 at index", i) // found 3 at index 2
	// }


	// Array Reverse - 2 ways

	// arr := [4]int{1, 2, 3, 4}
	// for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
	//     arr[i], arr[j] = arr[j], arr[i]
	// }
	// fmt.Println(arr) // [4 3 2 1]

	// 2nd way

	// arr := [4]int{1, 2, 3, 4}
	// sort.Sort(sort.Reverse(sort.IntSlice(arr[:])))
	// fmt.Println(arr) // [4 3 2 1]


	// Array of Slices -> slices is an array where each element is a slice
	// var arr [2][]int

	
	// Array of Maps -> array of maps is one in which each element is a map
	// var arr [2]map[string]int

	// arr[0] = map[string]int{"a": 1, "b": 2}
	// arr[1] = map[string]int{"c": 3, "d": 4}


	// Array of Functions -> An array of functions is one where each element is a function.

	// var arr [2]func(int) int

	// arr[0] = func(x int) int {
	//     return x * x
	// }
	// arr[1] = func(x int) int {
	//     return x * x * x
	// }


	// Array of Pointers 
	
	// var arr [2]*int

	// a := 1
	// b := 2
	// arr[0] = &a
	// arr[1] = &b

	// Array of Structs
	
	// type Person struct {
	//     Name string
	//     Age  int
	// }

	// var arr [2]Person

	// arr[0] = Person{Name: "John", Age: 20}
	// arr[1] = Person{Name: "Jane", Age: 21}


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

	//pointer to array
	var p *[5]int = &numbers
	fmt.Println("Pointer to numbers array:", p)
	
}
