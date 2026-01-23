package learning

import "fmt"

// ==========================================
// STRUCTS AND CUSTOM TYPES IN GO
// ==========================================
// Structs are typed collections of fields
// Custom types create new types from existing ones

// ========== BASIC STRUCT DEFINITION ==========

// Person represents a basic struct with various field types
type Person struct {
	FirstName string
	LastName  string
	Age       int
	Email     string
}

// Book demonstrates struct with different data types
type Book struct {
	Title       string
	Author      string
	Pages       int
	IsAvailable bool
	Price       float64
}

// DemonstrateBasicStructs shows struct creation and usage
func DemonstrateBasicStructs() {
	fmt.Println("\n========== BASIC STRUCTS ==========")

	// 1. Creating struct - Method 1: Field by field
	var person1 Person
	person1.FirstName = "John"
	person1.LastName = "Doe"
	person1.Age = 30
	person1.Email = "john@example.com"
	fmt.Printf("Person 1: %+v\n", person1) // %+v shows field names

	// 2. Creating struct - Method 2: Struct literal (all fields)
	person2 := Person{
		FirstName: "Jane",
		LastName:  "Smith",
		Age:       28,
		Email:     "jane@example.com",
	}
	fmt.Printf("Person 2: %+v\n", person2)

	// 3. Creating struct - Method 3: Short form (must match order)
	person3 := Person{"Bob", "Johnson", 35, "bob@example.com"}
	fmt.Printf("Person 3: %+v\n", person3)

	// 4. Creating struct - Method 4: Partial initialization
	person4 := Person{
		FirstName: "Alice",
		Age:       25,
		// Other fields get zero values: "" for string, 0 for int
	}
	fmt.Printf("Person 4: %+v\n", person4)

	// 5. Accessing and modifying struct fields
	fmt.Printf("\nPerson 2's full name: %s %s\n", person2.FirstName, person2.LastName)
	person2.Age = 29 // Modify field
	fmt.Printf("Person 2's updated age: %d\n", person2.Age)
}

// ========== NESTED STRUCTS ==========

// Address represents a nested struct
type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
}

// Employee demonstrates struct embedding/nesting
type Employee struct {
	ID      int
	Name    string
	Address Address // Nested struct
	Salary  float64
}

// DemonstrateNestedStructs shows nested struct usage
func DemonstrateNestedStructs() {
	fmt.Println("\n========== NESTED STRUCTS ==========")

	emp := Employee{
		ID:   101,
		Name: "Alice Cooper",
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			State:   "NY",
			ZipCode: "10001",
		},
		Salary: 75000.50,
	}

	fmt.Printf("Employee: %+v\n", emp)
	fmt.Printf("Employee lives in: %s, %s\n", emp.Address.City, emp.Address.State)
}

// ========== ANONYMOUS STRUCTS ==========

// DemonstrateAnonymousStructs shows anonymous struct usage
func DemonstrateAnonymousStructs() {
	fmt.Println("\n========== ANONYMOUS STRUCTS ==========")

	// Anonymous struct - useful for one-time use
	// Common in JSON handling, temporary data structures
	config := struct {
		Host string
		Port int
		SSL  bool
	}{
		Host: "localhost",
		Port: 8080,
		SSL:  false,
	}

	fmt.Printf("Config: %+v\n", config)
	fmt.Printf("Server: %s:%d, SSL: %v\n", config.Host, config.Port, config.SSL)
}

// ========== STRUCT COMPARISON ==========

// Point represents a 2D coordinate
type Point struct {
	X int
	Y int
}

// DemonstrateStructComparison shows struct comparison
func DemonstrateStructComparison() {
	fmt.Println("\n========== STRUCT COMPARISON ==========")

	point1 := Point{X: 10, Y: 20}
	point2 := Point{X: 10, Y: 20}
	point3 := Point{X: 15, Y: 25}

	// Structs can be compared if all fields are comparable
	fmt.Printf("point1 == point2: %v\n", point1 == point2) // true
	fmt.Printf("point1 == point3: %v\n", point1 == point3) // false

	// NOTE: Structs with slices, maps, or functions cannot be compared
	// because those types are not comparable
}

// ========== CUSTOM TYPES ==========

// Custom types create new types from existing ones
// Useful for type safety and adding methods

// UserID is a custom type based on int
type UserID int

// Status is a custom type based on string
type Status string

// Define constants for Status type (common pattern)
const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusPending  Status = "pending"
)

// DemonstrateCustomTypes shows custom type usage
func DemonstrateCustomTypes() {
	fmt.Println("\n========== CUSTOM TYPES ==========")

	// Custom types provide type safety
	var userID UserID = 12345
	var status Status = StatusActive

	fmt.Printf("User ID: %d (type: UserID)\n", userID)
	fmt.Printf("Status: %s (type: Status)\n", status)

	// Type safety - cannot directly assign between different custom types
	// even if underlying type is the same
	var normalInt int = 100
	// userID = normalInt // This would cause a compile error!
	userID = UserID(normalInt) // Must explicitly convert
	fmt.Printf("Converted User ID: %d\n", userID)

	// Custom types can have methods (see receiver_functions.go)
}

// ========== STRUCT WITH TAGS ==========

// User demonstrates struct tags (used with JSON, database ORMs, etc.)
type User struct {
	ID       int    `json:"id" db:"user_id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email_address"`
	IsActive bool   `json:"is_active" db:"active"`
}

// DemonstrateStructTags shows struct tag usage
func DemonstrateStructTags() {
	fmt.Println("\n========== STRUCT TAGS ==========")

	user := User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		IsActive: true,
	}

	fmt.Printf("User struct: %+v\n", user)
	fmt.Println("Struct tags are used by packages like encoding/json for serialization")
	fmt.Println("Tags define how fields map to JSON keys, database columns, etc.")
}

// ========== EMPTY STRUCT ==========

// DemonstrateEmptyStruct shows empty struct usage
func DemonstrateEmptyStruct() {
	fmt.Println("\n========== EMPTY STRUCT ==========")

	// Empty struct uses 0 bytes of memory
	// Useful for signals, sets, and when you only need existence
	type Signal struct{}

	signal := Signal{}
	fmt.Printf("Empty struct: %+v\n", signal)
	fmt.Println("Empty structs are often used in channels or as set values")

	// Example: Set using map with empty struct
	uniqueNumbers := make(map[int]struct{})
	uniqueNumbers[1] = struct{}{}
	uniqueNumbers[2] = struct{}{}
	uniqueNumbers[1] = struct{}{} // Duplicate, won't add

	fmt.Printf("Set contains: ")
	for num := range uniqueNumbers {
		fmt.Printf("%d ", num)
	}
	fmt.Println()
}

// ========== STRUCT POINTERS ==========

// DemonstrateStructPointers shows pointer usage with structs
func DemonstrateStructPointers() {
	fmt.Println("\n========== STRUCT POINTERS ==========")

	person := Person{
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
	}

	// Pointer to struct
	personPtr := &person

	// Go automatically dereferences pointers to access fields
	// Both syntaxes work:
	fmt.Printf("Age via pointer: %d\n", personPtr.Age)     // Automatic dereference
	fmt.Printf("Age via explicit: %d\n", (*personPtr).Age) // Explicit dereference

	// Modifying through pointer affects original
	personPtr.Age = 31
	fmt.Printf("Original person age: %d\n", person.Age) // Changed to 31

	// Creating struct with new (returns pointer)
	personPtr2 := new(Person)
	personPtr2.FirstName = "Jane"
	personPtr2.LastName = "Smith"
	personPtr2.Age = 28
	fmt.Printf("Person via new: %+v\n", personPtr2)
}

// RunAllStructDemos runs all struct demonstrations
func RunAllStructDemos() {
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║   LEARNING STRUCTS & CUSTOM TYPES     ║")
	fmt.Println("╚════════════════════════════════════════╝")

	DemonstrateBasicStructs()
	DemonstrateNestedStructs()
	DemonstrateAnonymousStructs()
	DemonstrateStructComparison()
	DemonstrateCustomTypes()
	DemonstrateStructTags()
	DemonstrateEmptyStruct()
	DemonstrateStructPointers()
}
