package learning

import "fmt"

// Interface definition
type Shape interface {
Area() float64
Perimeter() float64
}

// Circle type implementing Shape interface
type Circle struct {
Radius float64
}

func (c Circle) Area() float64 {
return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
return 2 * 3.14159 * c.Radius
}

// Square type implementing Shape interface
type Square struct {
Side float64
}

func (s Square) Area() float64 {
return s.Side * s.Side
}

func (s Square) Perimeter() float64 {
return 4 * s.Side
}

// Function that accepts interface
func printShapeInfo(s Shape) {
fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

// Example demonstrating interfaces
func ExampleInterfaces() {
circle := Circle{Radius: 5}
square := Square{Side: 4}

fmt.Println("Circle:")
printShapeInfo(circle)

fmt.Println("Square:")
printShapeInfo(square)
}
