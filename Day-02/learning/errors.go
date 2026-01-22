package learning

import "fmt"

//Errors in Go 
func ExampleErrors() {
_, err := Divide(10, 0)
if err != nil {
fmt.Printf("Error: %v\n", err)
}

result, err := Divide(10, 2)
if err != nil {
fmt.Printf("Error: %v\n", err)
} else {
fmt.Printf("10 / 2 = %d\n", result)
}

// Demonstrating error wrapping
_, err = Divide(5, 0)
if err != nil {
wrappedErr := fmt.Errorf("failed to perform division: %w", err)
fmt.Printf("Wrapped Error: %v\n", wrappedErr)
}

}

// Divide function that returns an error for division by zero
func Divide(a, b int) (int, error) {
if b == 0 {
return 0, fmt.Errorf("cannot divide by zero")
}
return a / b, nil
}
