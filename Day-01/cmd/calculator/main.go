package main

import (
	"fmt"
	"os"

	"day01/internal/calculator"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	a, b := 10, 5

	sum := calculator.Add(a, b)
	fmt.Printf("%d + %d = %d\n", a, b, sum)

	diff := calculator.Subtract(a, b)
	fmt.Printf("%d - %d = %d\n", a, b, diff)

	product := calculator.Multiply(a, b)
	fmt.Printf("%d * %d = %d\n", a, b, product)

	quotient, err := calculator.Divide(a, b)
	if err != nil {
		return fmt.Errorf("divide: %w", err)
	}
	fmt.Printf("%d / %d = %d\n", a, b, quotient)

	_, err = calculator.Divide(a, 0)
	if err != nil {
		fmt.Printf("%d / 0 = error: %v\n", a, err)
	}

	return nil
}
