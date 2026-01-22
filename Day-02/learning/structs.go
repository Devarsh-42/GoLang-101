package learning

import "fmt"

// Struct definition
type Person struct {
	Name string
	Age  int
}

// Function to demonstrate struct usage
func ExampleStructs() {
	p := Person{Name: "Alice", Age: 50}
	fmt.Printf("Name: %s, Age: %d\n", p.Name, p.Age)
}

// Struct with methods
type Rectangle struct {
	Width, Height float64
}

// Method to calculate area
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Function to demonstrate struct methods
func ExampleStructMethods() {
	rect := Rectangle{Width: 20, Height: 5}
	area := rect.Area()
	fmt.Printf("Area of rectangle: %.2f\n", area)
}
