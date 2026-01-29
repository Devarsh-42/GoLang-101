package main

import (
	"fmt"
	"math"
	"strings"
)

// =============================================================================
// INTERFACE DEFINITION
// =============================================================================

// Simple interface with one method
type Speaker interface {
	Speak() string
}

// Interface with multiple methods
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Empty interface - accepts any type
// interface{} or any (Go 1.18+) can hold values of any type
type AnyType interface{}

// =============================================================================
// IMPLEMENTING INTERFACES
// =============================================================================

// Dog implements Speaker interface
type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return "Woof! My name is " + d.Name
}

// Cat implements Speaker interface
type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return "Meow! My name is " + c.Name
}

// Human implements Speaker interface
type Human struct {
	Name string
}

func (h Human) Speak() string {
	return "Hello! My name is " + h.Name
}

// Rectangle implements Shape interface
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle implements Shape interface
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Triangle implements Shape interface
type Triangle struct {
	A, B, C float64 // Three sides
}

func (t Triangle) Area() float64 {
	// Using Heron's formula
	s := (t.A + t.B + t.C) / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

// =============================================================================
// INTERFACE AS FUNCTION PARAMETERS
// =============================================================================

// Function accepting interface - works with any type that implements Speaker
func introduce(s Speaker) {
	fmt.Println(s.Speak())
}

// Function accepting Shape interface
func printShapeInfo(s Shape) {
	fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

// Function accepting slice of interfaces
func introduceAll(speakers []Speaker) {
	for i, speaker := range speakers {
		fmt.Printf("%d. %s\n", i+1, speaker.Speak())
	}
}

// =============================================================================
// TYPE ASSERTION
// =============================================================================

// Type assertion - extracting concrete type from interface
func demonstrateTypeAssertion() {
	var s Speaker = Dog{Name: "Buddy"}

	// Type assertion - single value (panics if wrong type)
	// dog := s.(Dog)
	// fmt.Printf("Dog: %s\n", dog.Name)

	// Type assertion with ok check (safe, doesn't panic)
	if dog, ok := s.(Dog); ok {
		fmt.Printf("Successfully asserted as Dog: %s\n", dog.Name)
	} else {
		fmt.Println("Not a Dog")
	}

	// This will not panic, just return false
	if _, ok := s.(Cat); !ok {
		fmt.Println("Not a Cat")
	}
}

// Function using type assertion to get additional info
func getAnimalInfo(s Speaker) {
	fmt.Println(s.Speak())

	// Check specific type for type-specific behavior
	switch animal := s.(type) {
	case Dog:
		fmt.Printf("  This is a dog named %s\n", animal.Name)
	case Cat:
		fmt.Printf("  This is a cat named %s\n", animal.Name)
	case Human:
		fmt.Printf("  This is a human named %s\n", animal.Name)
	default:
		fmt.Printf("  Unknown type: %T\n", animal)
	}
}

// =============================================================================
// TYPE SWITCH
// =============================================================================

// Type switch - determining concrete type of interface
func describeType(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Integer: %d\n", v)
	case string:
		fmt.Printf("String: %s (length: %d)\n", v, len(v))
	case bool:
		fmt.Printf("Boolean: %t\n", v)
	case float64:
		fmt.Printf("Float: %.2f\n", v)
	case []int:
		fmt.Printf("Slice of ints: %v\n", v)
	case Speaker:
		fmt.Printf("Speaker: %s\n", v.Speak())
	case Shape:
		fmt.Printf("Shape with area: %.2f\n", v.Area())
	case nil:
		fmt.Println("Nil value")
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
}

// Processing different types
func processValue(val interface{}) interface{} {
	switch v := val.(type) {
	case int:
		return v * 2
	case string:
		return strings.ToUpper(v)
	case bool:
		return !v
	case float64:
		return v * 1.5
	default:
		return val
	}
}

// =============================================================================
// POLYMORPHISM
// =============================================================================

// Polymorphism - different types behaving differently through same interface
func demonstratePolymorphism() {
	// Different types, same interface
	speakers := []Speaker{
		Dog{Name: "Buddy"},
		Cat{Name: "Whiskers"},
		Human{Name: "Alice"},
	}

	fmt.Println("Polymorphism - all speak differently:")
	for _, speaker := range speakers {
		introduce(speaker)
	}
}

// Calculate total area of different shapes
func calculateTotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}

// =============================================================================
// INTERFACE COMPOSITION
// =============================================================================

// Composing interfaces - combining multiple interfaces
type Reader interface {
	Read() string
}

type Writer interface {
	Write(string)
}

// Composed interface
type ReadWriter interface {
	Reader
	Writer
}

// Implementation of Reader
type FileReader struct {
	Content string
}

func (f FileReader) Read() string {
	return f.Content
}

// Implementation of Writer
type FileWriter struct {
	Content string
}

func (f *FileWriter) Write(data string) {
	f.Content = data
}

// Implementation of ReadWriter (implements both)
type File struct {
	Content string
}

func (f *File) Read() string {
	return f.Content
}

func (f *File) Write(data string) {
	f.Content = data
}

// Using composed interface
func copyData(dst Writer, src Reader) {
	data := src.Read()
	dst.Write(data)
}

// =============================================================================
// ADVANCED INTERFACE PATTERNS
// =============================================================================

// StringerInterface (from fmt package)
// Types implementing this will have custom string representation
type StringerInterface interface {
	String() string
}

type Point struct {
	X, Y float64
}

func (p Point) String() string {
	return fmt.Sprintf("Point(%.2f, %.2f)", p.X, p.Y)
}

// Error interface - standard error handling
type CustomError struct {
	Code    int
	Message string
}

func (e CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// Comparable interface pattern
type Comparable interface {
	CompareTo(other Comparable) int
}

type Integer int

func (i Integer) CompareTo(other Comparable) int {
	otherInt := other.(Integer)
	if i < otherInt {
		return -1
	} else if i > otherInt {
		return 1
	}
	return 0
}

// Sort interface pattern (simplified)
type Sortable interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

type IntSlice []int

func (s IntSlice) Len() int           { return len(s) }
func (s IntSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s IntSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// =============================================================================
// INTERFACE WITH EMBEDDED STRUCTS
// =============================================================================

// Animal base type
type Animal struct {
	Name  string
	Sound string
	IsPet bool
}

// Pet interface
type Pet interface {
	Play() string
	Feed() string
}

// DogPet with embedded Animal
type DogPet struct {
	Animal
	Breed string
}

func (d DogPet) Speak() string {
	return fmt.Sprintf("%s says %s", d.Name, d.Sound)
}

func (d DogPet) Play() string {
	return fmt.Sprintf("%s is playing fetch!", d.Name)
}

func (d DogPet) Feed() string {
	return fmt.Sprintf("Feeding %s dog food", d.Name)
}

// CatPet with embedded Animal
type CatPet struct {
	Animal
	IndoorOnly bool
}

func (c CatPet) Speak() string {
	return fmt.Sprintf("%s says %s", c.Name, c.Sound)
}

func (c CatPet) Play() string {
	return fmt.Sprintf("%s is playing with yarn!", c.Name)
}

func (c CatPet) Feed() string {
	return fmt.Sprintf("Feeding %s cat food", c.Name)
}

// =============================================================================
// INTERFACE SATISFACTION AT COMPILE TIME
// =============================================================================

// Compile-time interface satisfaction check
// This ensures the type implements the interface
var _ Speaker = (*Dog)(nil)
var _ Shape = (*Circle)(nil)
var _ Pet = (*DogPet)(nil)

// =============================================================================
// INTERFACE WITH GENERICS (Go 1.18+)
// =============================================================================

// Container interface that can work with any type
type Container[T any] interface {
	Add(item T)
	Get(index int) T
	Size() int
}

// Simple implementation
type SliceContainer[T any] struct {
	items []T
}

func (s *SliceContainer[T]) Add(item T) {
	s.items = append(s.items, item)
}

func (s *SliceContainer[T]) Get(index int) T {
	return s.items[index]
}

func (s *SliceContainer[T]) Size() int {
	return len(s.items)
}

// =============================================================================
// EMPTY INTERFACE USAGE
// =============================================================================

// Function accepting any type
func printAnything(val interface{}) {
	fmt.Printf("Value: %v, Type: %T\n", val, val)
}

// Map with any value type
func createFlexibleMap() map[string]interface{} {
	return map[string]interface{}{
		"name":    "Alice",
		"age":     30,
		"active":  true,
		"balance": 1234.56,
		"tags":    []string{"premium", "verified"},
	}
}

// =============================================================================
// DEMONSTRATION FUNCTION
// =============================================================================

// DemonstrateInterfaces showcases all interface concepts
func DemonstrateInterfaces() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("INTERFACES DEMONSTRATION")
	fmt.Println(strings.Repeat("=", 80))

	// 1. Basic Interface Implementation
	fmt.Println("\n1. BASIC INTERFACE IMPLEMENTATION:")

	dog := Dog{Name: "Buddy"}
	cat := Cat{Name: "Whiskers"}
	human := Human{Name: "Alice"}

	introduce(dog)
	introduce(cat)
	introduce(human)

	// 2. Shape Interface
	fmt.Println("\n2. SHAPE INTERFACE:")

	rect := Rectangle{Width: 5, Height: 3}
	circle := Circle{Radius: 4}
	triangle := Triangle{A: 3, B: 4, C: 5}

	fmt.Print("Rectangle: ")
	printShapeInfo(rect)
	fmt.Print("Circle: ")
	printShapeInfo(circle)
	fmt.Print("Triangle: ")
	printShapeInfo(triangle)

	// 3. Polymorphism
	fmt.Println("\n3. POLYMORPHISM:")
	demonstratePolymorphism()

	shapes := []Shape{rect, circle, triangle}
	totalArea := calculateTotalArea(shapes)
	fmt.Printf("Total area of all shapes: %.2f\n", totalArea)

	// 4. Type Assertion
	fmt.Println("\n4. TYPE ASSERTION:")
	demonstrateTypeAssertion()

	speakers := []Speaker{dog, cat, human}
	for _, speaker := range speakers {
		getAnimalInfo(speaker)
		fmt.Println()
	}

	// 5. Type Switch
	fmt.Println("\n5. TYPE SWITCH:")

	values := []interface{}{
		42,
		"Hello, World",
		true,
		3.14159,
		[]int{1, 2, 3},
		dog,
		circle,
		nil,
	}

	for _, val := range values {
		describeType(val)
	}

	// Process values
	fmt.Println("\nProcessing values:")
	for _, val := range []interface{}{10, "hello", true, 2.5} {
		result := processValue(val)
		fmt.Printf("%v (%T) -> %v (%T)\n", val, val, result, result)
	}

	// 6. Interface Composition
	fmt.Println("\n6. INTERFACE COMPOSITION:")

	file := &File{Content: "Initial content"}
	fmt.Printf("Read: %s\n", file.Read())

	file.Write("New content")
	fmt.Printf("After write: %s\n", file.Read())

	// Copy using composed interface
	src := FileReader{Content: "Source data"}
	dst := &FileWriter{}
	copyData(dst, src)
	fmt.Printf("Copied data: %s\n", dst.Content)

	// 7. Custom String Representation
	fmt.Println("\n7. STRINGER INTERFACE:")

	p := Point{X: 3.5, Y: 7.2}
	fmt.Println(p) // Uses String() method automatically

	// 8. Error Interface
	fmt.Println("\n8. ERROR INTERFACE:")

	err := CustomError{Code: 404, Message: "Not Found"}
	fmt.Printf("Error: %v\n", err)

	// 9. Embedded Structs with Interfaces
	fmt.Println("\n9. EMBEDDED STRUCTS WITH INTERFACES:")

	dogPet := DogPet{
		Animal: Animal{Name: "Max", Sound: "Woof", IsPet: true},
		Breed:  "Labrador",
	}

	catPet := CatPet{
		Animal:     Animal{Name: "Luna", Sound: "Meow", IsPet: true},
		IndoorOnly: true,
	}

	pets := []Pet{dogPet, catPet}
	for _, pet := range pets {
		if speaker, ok := pet.(Speaker); ok {
			fmt.Println(speaker.Speak())
		}
		fmt.Println(pet.Play())
		fmt.Println(pet.Feed())
		fmt.Println()
	}

	// 10. Generic Interfaces
	fmt.Println("\n10. GENERIC INTERFACES:")

	intContainer := &SliceContainer[int]{}
	intContainer.Add(10)
	intContainer.Add(20)
	intContainer.Add(30)

	fmt.Printf("Container size: %d\n", intContainer.Size())
	fmt.Printf("First item: %d\n", intContainer.Get(0))

	stringContainer := &SliceContainer[string]{}
	stringContainer.Add("Hello")
	stringContainer.Add("World")

	fmt.Printf("String container size: %d\n", stringContainer.Size())
	fmt.Printf("First item: %s\n", stringContainer.Get(0))

	// 11. Empty Interface
	fmt.Println("\n11. EMPTY INTERFACE (interface{}):")

	printAnything(42)
	printAnything("Hello")
	printAnything(true)
	printAnything([]int{1, 2, 3})

	flexMap := createFlexibleMap()
	fmt.Println("\nFlexible map:")
	for key, value := range flexMap {
		fmt.Printf("  %s: %v (%T)\n", key, value, value)
	}

	// 12. Interface Comparison
	fmt.Println("\n12. INTERFACE COMPARISON:")

	var s1 Speaker = Dog{Name: "Buddy"}
	var s2 Speaker = Dog{Name: "Buddy"}
	var s3 Speaker = Cat{Name: "Whiskers"}

	// Can compare interfaces if underlying types are comparable
	fmt.Printf("s1 == s2: %t\n", s1 == s2)
	fmt.Printf("s1 == s3: %t\n", s1 == s3)

	// 13. Nil Interface
	fmt.Println("\n13. NIL INTERFACE:")

	var nilSpeaker Speaker
	fmt.Printf("Nil speaker: %v\n", nilSpeaker)
	fmt.Printf("Is nil: %t\n", nilSpeaker == nil)

	// Be careful - interface with nil concrete value is not nil!
	var nilDog *Dog
	var speakerWithNilValue Speaker = nilDog
	fmt.Printf("Speaker with nil value == nil: %t\n", speakerWithNilValue == nil) // false!

	fmt.Println("\n" + strings.Repeat("=", 80))
}
