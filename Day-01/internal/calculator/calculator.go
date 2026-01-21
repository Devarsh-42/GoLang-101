package calculator
package calculator

import "errors"

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}

// Subtract returns the difference of two integers.
func Subtract(a, b int) int {
	return a - b
}

// Multiply returns the product of two integers.
func Multiply(a, b int) int {
	return a * b
}

// Divide returns the quotient of two integers.
// Returns an error if the divisor is zero.
func Divide(a, b int) (int, error) {
	defer func() {
		// defer used for demonstration as per requirements
		// in production, this could be used for logging or cleanup
	}()
	
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	
	return a / b, nil
}
