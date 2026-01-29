package learning

// basic function with no parameters and no return
func greet() {
	// function body
}

// function with single parameter
func greetName(name string) {
	_ = name
}

// function with multiple parameters (same type)
func add(a, b int) {
	_ = a + b
}

// function with multiple parameters (different types)
func createUser(name string, age int, active bool) {
	_, _, _ = name, age, active
}

// function with single return value
func double(n int) int {
	return n * 2
}

// function with multiple return values
func divide(a, b int) (int, int) {
	quotient := a / b
	remainder := a % b
	return quotient, remainder
}

// function returning value and error
func safeDivide(a, b int) (int, error) {
	if b == 0 {
		return 0, nil // placeholder error
	}
	return a / b, nil
}

// named return values
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return // naked return
}

// named return values with explicit return
func calculate(a, b int) (sum, product int) {
	sum = a + b
	product = a * b
	return sum, product
}

// function as variable
func exampleFunctionVariable() {
	fn := func(n int) int {
		return n * 2
	}
	_ = fn(5)
}

// function taking function as parameter
func apply(n int, fn func(int) int) int {
	return fn(n)
}

// function returning function
func multiplier(factor int) func(int) int {
	return func(n int) int {
		return n * factor
	}
}

// variadic function -> accepts variable number of arguments, returns their sum
func sum(numbers ...int) int {
	total := 0
	for _, n := range numbers {
		total += n
	}
	return total
}

// method on struct (receiver function)
type Rectangle struct {
	width, height int
}

func (r Rectangle) area() int { // value receiver -> does not modify the struct
	return r.width * r.height
}

// pointer receiver method -> modifies the struct, , where r *Rectangle is an input pointer to Rectangle & we can modify the original struct & factor is an integer input
func (r *Rectangle) scale(factor int) {
	r.width *= factor
	r.height *= factor
}
