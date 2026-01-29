package main

import (
	"fmt"
	"strings"
)

// =============================================================================
// BASIC FUNCTIONS
// =============================================================================

// Simple function with parameters and return value
func add(a, b int) int {
	return a + b
}

// Function with multiple parameters of same type (shorthand)
func multiply(x, y, z int) int {
	return x * y * z
}

// Function with no return value
func printMessage(message string) {
	fmt.Println(message)
}

// =============================================================================
// MULTIPLE RETURN VALUES
// =============================================================================

// Function returning multiple values
// Common pattern in Go for returning result and error
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// Function returning multiple values of same type
func swap(a, b string) (string, string) {
	return b, a
}

// Function returning multiple values with different types
func getUserInfo(id int) (string, int, bool) {
	// Simulating database lookup
	if id == 1 {
		return "Alice", 30, true // name, age, found
	}
	return "", 0, false
}

// =============================================================================
// NAMED RETURN VALUES
// =============================================================================

// Named return values - automatically initialized to zero values
// Can use naked return statement
func rectangle(length, width float64) (area, perimeter float64) {
	area = length * width
	perimeter = 2 * (length + width)
	return // naked return - returns named values
}

// Named return with explicit return values
func calculate(a, b int) (sum, product, diff int) {
	sum = a + b
	product = a * b
	diff = a - b
	return sum, product, diff // explicit return
}

// =============================================================================
// VARIADIC FUNCTIONS
// =============================================================================

// Variadic function - accepts variable number of arguments
// The variadic parameter must be the last parameter
func sum(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

// Variadic function with other parameters
// Regular parameters must come before variadic parameter
func formatMessage(prefix string, messages ...string) string {
	if len(messages) == 0 {
		return prefix
	}
	return prefix + ": " + strings.Join(messages, ", ")
}

// Variadic function with different operations
func processNumbers(operation string, numbers ...float64) float64 {
	if len(numbers) == 0 {
		return 0
	}

	switch operation {
	case "sum":
		total := 0.0
		for _, num := range numbers {
			total += num
		}
		return total
	case "avg":
		total := 0.0
		for _, num := range numbers {
			total += num
		}
		return total / float64(len(numbers))
	case "max":
		max := numbers[0]
		for _, num := range numbers {
			if num > max {
				max = num
			}
		}
		return max
	default:
		return 0
	}
}

// =============================================================================
// ANONYMOUS FUNCTIONS
// =============================================================================

// Function that demonstrates anonymous functions
func demonstrateAnonymousFunctions() {
	// Anonymous function assigned to a variable
	greet := func(name string) string {
		return "Hello, " + name
	}
	fmt.Println(greet("Bob"))

	// Immediately invoked anonymous function
	result := func(a, b int) int {
		return a * b
	}(5, 3) // Called immediately with arguments
	fmt.Println("Immediate result:", result)

	// Anonymous function with multiple statements
	calculate := func(x, y int) (int, int) {
		sum := x + y
		product := x * y
		return sum, product
	}
	s, p := calculate(4, 5)
	fmt.Printf("Sum: %d, Product: %d\n", s, p)
}

// =============================================================================
// CLOSURES
// =============================================================================

// Closure - function that references variables from outside its body
// The function "closes over" the variables it references
func counter() func() int {
	count := 0 // This variable is captured by the closure
	return func() int {
		count++ // Each call increments the same count variable
		return count
	}
}

// Closure that captures and modifies external state
func makeMultiplier(factor int) func(int) int {
	return func(x int) int {
		return x * factor
	}
}

// Closure with multiple captured variables
func makeAccumulator() func(int) int {
	sum := 0
	return func(value int) int {
		sum += value
		return sum
	}
}

// Advanced closure - returns multiple functions sharing same state
func makeCounter() (increment func() int, decrement func() int, reset func()) {
	count := 0

	increment = func() int {
		count++
		return count
	}

	decrement = func() int {
		count--
		return count
	}

	reset = func() {
		count = 0
	}

	return
}

// =============================================================================
// HIGHER-ORDER FUNCTIONS
// =============================================================================

// Higher-order function - function that takes another function as parameter
func applyOperation(a, b int, operation func(int, int) int) int {
	return operation(a, b)
}

// Higher-order function that returns a function
func makeAdder(x int) func(int) int {
	return func(y int) int {
		return x + y
	}
}

// Function that takes a function and applies it to a slice
func mapInts(numbers []int, fn func(int) int) []int {
	result := make([]int, len(numbers))
	for i, num := range numbers {
		result[i] = fn(num)
	}
	return result
}

// Filter function - returns slice with elements that satisfy predicate
func filterInts(numbers []int, predicate func(int) bool) []int {
	result := []int{}
	for _, num := range numbers {
		if predicate(num) {
			result = append(result, num)
		}
	}
	return result
}

// Reduce function - reduces slice to single value
func reduceInts(numbers []int, initial int, reducer func(int, int) int) int {
	result := initial
	for _, num := range numbers {
		result = reducer(result, num)
	}
	return result
}

// Function composition - combines multiple functions
func compose(f, g func(int) int) func(int) int {
	return func(x int) int {
		return f(g(x))
	}
}

// =============================================================================
// CALL BY VALUE
// =============================================================================

// Go is always call by value - it passes copies of values
func modifyValue(x int) {
	x = 100 // This only modifies the local copy
}

func modifySlice(s []int) {
	// Slices are reference types - the slice header is copied,
	// but it still points to the same underlying array
	s[0] = 999 // This WILL modify the original array
}

func modifySliceStructure(s []int) {
	// Appending creates a new underlying array if capacity is exceeded
	// This won't affect the original slice
	s = append(s, 100)
}

// Using pointers to modify values
func modifyWithPointer(x *int) {
	*x = 200 // Dereference and modify the value at the pointer address
}

// =============================================================================
// RECURSIVE FUNCTIONS
// =============================================================================

// Factorial using recursion
func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

// Fibonacci using recursion
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// Fibonacci with memoization (closure + recursion)
func fibonacciMemo() func(int) int {
	cache := make(map[int]int)

	var fib func(int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}

		// Check cache first
		if val, ok := cache[n]; ok {
			return val
		}

		// Calculate and cache
		result := fib(n-1) + fib(n-2)
		cache[n] = result
		return result
	}

	return fib
}

// =============================================================================
// FUNCTION AS TYPE
// =============================================================================

// Define custom function type
type MathOperation func(int, int) int

// Function that uses custom function type
func executeOperation(a, b int, op MathOperation) int {
	return op(a, b)
}

// Define a struct with function fields
type Calculator struct {
	Add      func(int, int) int
	Subtract func(int, int) int
	Multiply func(int, int) int
}

func newCalculator() Calculator {
	return Calculator{
		Add:      func(a, b int) int { return a + b },
		Subtract: func(a, b int) int { return a - b },
		Multiply: func(a, b int) int { return a * b },
	}
}

// =============================================================================
// DEFER, PANIC, AND RECOVER
// =============================================================================

// Defer - schedules function call to be run after the function completes
func demonstrateDefer() {
	defer fmt.Println("This runs last")
	defer fmt.Println("This runs second")
	fmt.Println("This runs first")
	// Output order: first, second, last (LIFO - Last In First Out)
}

// Practical defer usage - cleanup operations
func processFile(filename string) error {
	fmt.Printf("Opening file: %s\n", filename)
	// In real code: file, err := os.Open(filename)
	defer fmt.Println("Closing file") // Ensures cleanup

	fmt.Println("Processing file...")
	return nil
}

// =============================================================================
// DEMONSTRATION FUNCTION
// =============================================================================

// DemonstrateFunctions showcases all function concepts
func DemonstrateFunctions() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("FUNCTIONS DEMONSTRATION")
	fmt.Println(strings.Repeat("=", 80))

	// 1. Basic Functions
	fmt.Println("\n1. BASIC FUNCTIONS:")
	fmt.Printf("add(5, 3) = %d\n", add(5, 3))
	fmt.Printf("multiply(2, 3, 4) = %d\n", multiply(2, 3, 4))
	printMessage("Hello from basic function!")

	// 2. Multiple Return Values
	fmt.Println("\n2. MULTIPLE RETURN VALUES:")
	result, err := divide(10, 2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("divide(10, 2) = %.2f\n", result)
	}

	_, err = divide(10, 0)
	fmt.Printf("divide(10, 0) error: %v\n", err)

	a, b := swap("first", "second")
	fmt.Printf("swap('first', 'second') = '%s', '%s'\n", a, b)

	// 3. Named Return Values
	fmt.Println("\n3. NAMED RETURN VALUES:")
	area, perimeter := rectangle(5, 3)
	fmt.Printf("rectangle(5, 3) - Area: %.2f, Perimeter: %.2f\n", area, perimeter)

	// 4. Variadic Functions
	fmt.Println("\n4. VARIADIC FUNCTIONS:")
	fmt.Printf("sum(1, 2, 3, 4, 5) = %d\n", sum(1, 2, 3, 4, 5))
	fmt.Printf("sum(10, 20) = %d\n", sum(10, 20))

	numbers := []int{1, 2, 3, 4, 5}
	fmt.Printf("sum(numbers...) = %d\n", sum(numbers...)) // Spread operator

	msg := formatMessage("ERROR", "File not found", "Invalid permissions")
	fmt.Println(msg)

	fmt.Printf("processNumbers('avg', 10, 20, 30, 40) = %.2f\n",
		processNumbers("avg", 10, 20, 30, 40))

	// 5. Anonymous Functions
	fmt.Println("\n5. ANONYMOUS FUNCTIONS:")
	demonstrateAnonymousFunctions()

	// 6. Closures
	fmt.Println("\n6. CLOSURES:")

	// Simple counter closure
	cnt := counter()
	fmt.Printf("Counter: %d\n", cnt())
	fmt.Printf("Counter: %d\n", cnt())
	fmt.Printf("Counter: %d\n", cnt())

	// Multiple closures with independent state
	cnt1 := counter()
	cnt2 := counter()
	fmt.Printf("Counter1: %d, Counter2: %d\n", cnt1(), cnt2())
	fmt.Printf("Counter1: %d, Counter2: %d\n", cnt1(), cnt2())

	// Closure capturing factor
	double := makeMultiplier(2)
	triple := makeMultiplier(3)
	fmt.Printf("double(5) = %d, triple(5) = %d\n", double(5), triple(5))

	// Accumulator closure
	acc := makeAccumulator()
	fmt.Printf("Accumulator: %d\n", acc(5))
	fmt.Printf("Accumulator: %d\n", acc(10))
	fmt.Printf("Accumulator: %d\n", acc(15))

	// Advanced closure with multiple functions
	inc, dec, reset := makeCounter()
	fmt.Printf("Increment: %d\n", inc())
	fmt.Printf("Increment: %d\n", inc())
	fmt.Printf("Decrement: %d\n", dec())
	reset()
	fmt.Printf("After reset, Increment: %d\n", inc())

	// 7. Higher-order Functions
	fmt.Println("\n7. HIGHER-ORDER FUNCTIONS:")

	addOp := func(x, y int) int { return x + y }
	mulOp := func(x, y int) int { return x * y }

	fmt.Printf("applyOperation(5, 3, add) = %d\n", applyOperation(5, 3, addOp))
	fmt.Printf("applyOperation(5, 3, multiply) = %d\n", applyOperation(5, 3, mulOp))

	add10 := makeAdder(10)
	fmt.Printf("add10(5) = %d\n", add10(5))

	// Map, Filter, Reduce
	nums := []int{1, 2, 3, 4, 5}

	squared := mapInts(nums, func(x int) int { return x * x })
	fmt.Printf("Map (square): %v -> %v\n", nums, squared)

	evens := filterInts(nums, func(x int) bool { return x%2 == 0 })
	fmt.Printf("Filter (evens): %v -> %v\n", nums, evens)

	sumResult := reduceInts(nums, 0, func(acc, x int) int { return acc + x })
	fmt.Printf("Reduce (sum): %v -> %d\n", nums, sumResult)

	// Function composition
	addOne := func(x int) int { return x + 1 }
	multiplyByTwo := func(x int) int { return x * 2 }
	addThenMultiply := compose(multiplyByTwo, addOne)
	fmt.Printf("compose(×2, +1)(5) = %d\n", addThenMultiply(5)) // (5+1)*2 = 12

	// 8. Call by Value
	fmt.Println("\n8. CALL BY VALUE:")

	x := 50
	fmt.Printf("Before modifyValue: x = %d\n", x)
	modifyValue(x)
	fmt.Printf("After modifyValue: x = %d (unchanged)\n", x)

	slice := []int{1, 2, 3}
	fmt.Printf("Before modifySlice: %v\n", slice)
	modifySlice(slice)
	fmt.Printf("After modifySlice: %v (modified)\n", slice)

	fmt.Printf("Before modifyWithPointer: x = %d\n", x)
	modifyWithPointer(&x)
	fmt.Printf("After modifyWithPointer: x = %d (modified)\n", x)

	// 9. Recursive Functions
	fmt.Println("\n9. RECURSIVE FUNCTIONS:")
	fmt.Printf("factorial(5) = %d\n", factorial(5))
	fmt.Printf("fibonacci(10) = %d\n", fibonacci(10))

	// Memoized fibonacci
	fib := fibonacciMemo()
	fmt.Printf("fibonacciMemo(20) = %d\n", fib(20))
	fmt.Printf("fibonacciMemo(30) = %d (cached)\n", fib(30))

	// 10. Function Types
	fmt.Println("\n10. FUNCTION TYPES:")

	var op MathOperation = func(a, b int) int { return a + b }
	fmt.Printf("Custom type operation: %d\n", executeOperation(5, 3, op))

	calc := newCalculator()
	fmt.Printf("Calculator - Add: %d, Subtract: %d, Multiply: %d\n",
		calc.Add(10, 5), calc.Subtract(10, 5), calc.Multiply(10, 5))

	// 11. Defer
	fmt.Println("\n11. DEFER:")
	demonstrateDefer()
	processFile("example.txt")

	fmt.Println("\n" + strings.Repeat("=", 80))
}
