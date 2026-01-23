# Day 03 - Maps, Structs, Custom Types & Receiver Functions

This directory contains comprehensive learning materials for Day 3 concepts in Go.

## Project Structure

```
Day-03/
├── go.mod                              # Go module file
├── README.md                           # This file
├── cmd/
│   └── learning-demo/
│       └── main.go                     # Interactive learning demo
└── learning/
    ├── maps.go                         # Map concepts and examples
    ├── struct_and_custome_types.go     # Structs and custom types
    └── reciver_functions.go            # Receiver functions (methods)
```

## How to Run

From the Day-03 directory, run:

```bash
go run ./cmd/learning-demo/main.go
```

## Topics Covered

### 1. Maps (`learning/maps.go`)
- **Map Basics**: Creating, reading, updating, and deleting
- **Map Iteration**: Different ways to iterate over maps
- **Complex Maps**: Nested maps, maps with structs as values
- **Reference Types**: Understanding how maps are passed
- **Practical Example**: Word counter implementation

#### Key Concepts:
```go
// Creating maps
m := make(map[string]int)
m := map[string]int{"key": value}

// Operations
m[key] = value              // Add/Update
value, ok := m[key]         // Get with existence check
delete(m, key)              // Delete
len(m)                      // Size

// Iteration
for key, value := range m {
    // Process key-value pairs
}
```

### 2. Structs & Custom Types (`learning/struct_and_custome_types.go`)
- **Basic Structs**: Declaration and initialization
- **Nested Structs**: Embedding structs within structs
- **Anonymous Structs**: One-time use structs
- **Struct Comparison**: When and how structs can be compared
- **Custom Types**: Creating new types from existing ones
- **Struct Tags**: Metadata for serialization
- **Empty Structs**: Zero-memory structs for signals
- **Struct Pointers**: Working with pointers to structs

#### Key Concepts:
```go
// Define struct
type Person struct {
    Name string
    Age  int
}

// Create struct
p1 := Person{Name: "John", Age: 30}
p2 := Person{"Jane", 25}
p3 := &Person{Name: "Bob"}  // Pointer

// Custom types
type UserID int
type Status string

// Type safety
var id UserID = 12345
```

### 3. Receiver Functions (`learning/reciver_functions.go`)
- **Value Receivers**: Methods that work on copies
- **Pointer Receivers**: Methods that modify the original
- **Methods on Custom Types**: Attaching behavior to custom types
- **Methods with Parameters**: Complex method signatures
- **Methods vs Functions**: Understanding the differences
- **Method Chaining**: Creating fluent APIs

#### Key Concepts:
```go
// Value receiver (doesn't modify original)
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Pointer receiver (modifies original)
func (c *Circle) Scale(factor float64) {
    c.Radius *= factor
}

// Usage
rect := Rectangle{Width: 10, Height: 5}
area := rect.Area()
rect.Scale(2)  // Won't work for value receivers
```

## Learning Approach

Each file in the `learning/` directory contains:
1. **Clear Comments**: Explaining what each concept does
2. **Multiple Examples**: Different use cases and patterns
3. **Practical Code**: Real-world applicable examples
4. **Key Takeaways**: Summary of important points

## Interactive Menu

The main program provides an interactive menu where you can:
- Learn about each topic individually
- Run all demonstrations at once
- See practical examples with output
- Revise concepts at your own pace

## Tips for Revision

1. **Read the Code**: Each function is well-documented with comments
2. **Run Examples**: Try option 1-3 to see individual topics
3. **Modify & Experiment**: Change values and see what happens
4. **Use Option 4**: Run all demos to see complete flow
5. **Check Comments**: Quick reference guide at bottom of main.go

## Quick Reference

### Maps
- Use `make(map[K]V)` to create
- Always check existence: `value, ok := m[key]`
- Iteration order is NOT guaranteed
- Maps are reference types

### Structs
- Group related data together
- Can be compared if all fields are comparable
- Use struct tags for serialization
- Prefer pointers for large structs

### Receiver Functions
- Value receiver: `func (r Type) Method()`
- Pointer receiver: `func (r *Type) Method()`
- Use pointer receivers when:
  - Method needs to modify receiver
  - Receiver is large
  - For consistency across methods

## Key Differences

| Feature | Value Receiver | Pointer Receiver |
|---------|---------------|------------------|
| Modifies original | No | Yes |
| Memory usage | Higher (copy) | Lower (reference) |
| Use when | Small types, read-only | Large types, mutations |
| Syntax | `func (t Type)` | `func (t *Type)` |

## Best Practices

1. **Maps**: Always check if key exists before using value
2. **Structs**: Use meaningful field names, consider tags for JSON
3. **Methods**: Be consistent - if one uses pointer, all should
4. **Custom Types**: Use for type safety and domain modeling
5. **Comments**: Document exported types and functions

## Examples to Try

After understanding the code, try these exercises:
1. Create a student management system using maps and structs
2. Implement a temperature converter with custom types and methods
3. Build a bank account system with receiver functions
4. Create a shape library with area/perimeter calculations

---

**Happy Learning!**

To revise any concept, simply:
1. Run the program
2. Choose the topic number
3. Read the output carefully
4. Review the corresponding source file
5. Experiment with modifications
