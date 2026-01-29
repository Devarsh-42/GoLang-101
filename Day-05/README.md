# Day 05 - Advanced Go Concepts

This directory contains comprehensive examples and demonstrations of advanced Go programming concepts including **Functions**, **Maps**, **Interfaces**, and **Generics**.

## 📁 File Structure

```
Day-05/
├── main.go          # Main entry point that runs all demonstrations
├── functions.go     # Comprehensive function concepts and patterns
├── maps.go          # Complete map operations and use cases
├── interfaces.go    # Interface definitions, polymorphism, and composition
├── generics.go      # Generic programming with type parameters
├── go.mod           # Module definition with dependencies
└── README.md        # This file
```

## 🚀 Running the Code

To run all demonstrations:

```bash
cd Day-05
go run main.go functions.go maps.go interfaces.go generics.go
```

Or build and run:

```bash
go build
./day05
```

## 📚 Topics Covered

### 1. **functions.go** - Functions in Go

#### Basic Concepts
- **Function Basics**: Simple functions with parameters and return values
- **Multiple Return Values**: Returning multiple values (common for result + error)
- **Named Return Values**: Pre-declared return values with naked return
- **Variadic Functions**: Functions accepting variable number of arguments
- **Function as Type**: Custom function types and function fields in structs

#### Advanced Patterns
- **Anonymous Functions**: Function literals and immediate invocation
- **Closures**: Functions that capture external variables
- **Higher-order Functions**: Functions that accept or return other functions
- **Function Composition**: Combining multiple functions
- **Defer Statements**: Scheduling function calls for cleanup

#### Key Examples
```go
// Variadic function
func sum(numbers ...int) int

// Closure
func counter() func() int

// Higher-order function
func applyOperation(a, b int, operation func(int, int) int) int

// Map/Filter/Reduce patterns
func mapInts(numbers []int, fn func(int) int) []int
func filterInts(numbers []int, predicate func(int) bool) []int
func reduceInts(numbers []int, initial int, reducer func(int, int) int) int
```

#### Important Concepts
- ✅ Go is always **call by value** (passes copies)
- ✅ Closures can capture and modify external state
- ✅ Defer statements execute in **LIFO** (Last In First Out) order
- ✅ Multiple return values are idiomatic for error handling

---

### 2. **maps.go** - Maps (Hash Tables)

#### Creation Methods
- Using `make()` - preferred method
- Map literals with initialization
- Empty map literals
- Nil maps (read-safe, write panics)

#### Basic Operations
- **Adding/Updating**: `map[key] = value`
- **Deleting**: `delete(map, key)`
- **Checking existence**: `value, exists := map[key]`
- **Getting length**: `len(map)`

#### Advanced Topics
- **Iteration**: Unordered iteration, sorted iteration
- **Complex Types**: Maps with struct values, slice values, map values, struct keys
- **Reference Semantics**: Maps are reference types (mutations are visible)
- **Advanced Patterns**:
  - Map as Set (using `map[T]struct{}`)
  - Counting occurrences
  - Grouping items
  - Inverting maps
  - Merging maps
  - Default values pattern

#### Key Examples
```go
// Map as set
set := make(map[string]struct{})
set["item"] = struct{}{}

// Counting occurrences
counts := make(map[string]int)
for _, word := range words {
    counts[word]++ // Auto-initializes to 0
}

// Nested maps
config := map[string]map[string]string{
    "database": {"host": "localhost", "port": "5432"},
}
```

#### Important Concepts
- ✅ Maps are **unordered** (iteration order is random)
- ✅ Maps are **reference types** (assignment creates reference, not copy)
- ✅ Reading from nil map returns zero value (safe)
- ✅ Writing to nil map **panics**
- ✅ Map keys must be **comparable** (can use ==)

---

### 3. **interfaces.go** - Interfaces and Polymorphism

#### Core Concepts
- **Interface Definition**: Sets of method signatures
- **Implicit Implementation**: No "implements" keyword needed
- **Type Assertion**: Extracting concrete type from interface
- **Type Switch**: Determining concrete type at runtime
- **Empty Interface**: `interface{}` accepts any type

#### Advanced Topics
- **Polymorphism**: Different types with same interface behaving differently
- **Interface Composition**: Combining multiple interfaces
- **Embedded Interfaces**: Building complex interfaces from simple ones
- **Custom Error Types**: Implementing `error` interface
- **Stringer Interface**: Custom string representation

#### Key Examples
```go
// Interface definition
type Speaker interface {
    Speak() string
}

// Implicit implementation
type Dog struct { Name string }
func (d Dog) Speak() string { return "Woof!" }

// Type assertion
if dog, ok := speaker.(Dog); ok {
    fmt.Println(dog.Name)
}

// Type switch
switch v := i.(type) {
case int:
    fmt.Printf("Integer: %d\n", v)
case string:
    fmt.Printf("String: %s\n", v)
}

// Interface composition
type ReadWriter interface {
    Reader
    Writer
}
```

#### Important Concepts
- ✅ Interfaces enable **polymorphism** without inheritance
- ✅ Implementation is **implicit** (duck typing)
- ✅ Empty interface can hold **any value**
- ✅ Interface with nil concrete value is **not nil**
- ✅ Compile-time check: `var _ Interface = (*Type)(nil)`

---

### 4. **generics.go** - Generic Programming

#### Type Parameters
- **Basic Generics**: Type parameters with `any` constraint
- **Multiple Type Parameters**: Functions with multiple generic types
- **Generic Types**: Structs and types with type parameters

#### Constraints
- **`any`**: No restrictions (accepts any type)
- **`comparable`**: Types supporting `==` and `!=`
- **`constraints.Ordered`**: Types supporting `<`, `>`, `<=`, `>=`
- **Custom Constraints**: Define your own using interfaces

#### Generic Data Structures
- **Stack**: Generic LIFO structure
- **Queue**: Generic FIFO structure
- **LinkedList**: Generic linked list
- **Optional**: Type-safe optional values (like Rust's Option)
- **Result**: Type-safe error handling (like Rust's Result)

#### Functional Programming
- **Map**: Transform slice elements
- **Filter**: Select elements matching predicate
- **Reduce**: Reduce slice to single value
- **FlatMap**: Map and flatten results
- **GroupBy**: Group elements by key
- **Partition**: Split slice by predicate
- **All/Any**: Check if all/any elements satisfy predicate

#### Key Examples
```go
// Generic function
func Print[T any](value T)

// Comparable constraint
func Contains[T comparable](slice []T, value T) bool

// Ordered constraint
func Min[T constraints.Ordered](a, b T) T

// Custom constraint
type Number interface {
    ~int | ~int8 | ... | ~float64
}

// Generic stack
type Stack[T any] struct {
    items []T
}

// Map/Filter/Reduce
squared := Map(nums, func(x int) int { return x * x })
evens := Filter(nums, func(x int) bool { return x%2 == 0 })
sum := Reduce(nums, 0, func(acc, x int) int { return acc + x })
```

#### Important Concepts
- ✅ Generics enable **type-safe** reusable code
- ✅ Constraints restrict what operations can be performed
- ✅ Use `~` in constraints for **underlying types**
- ✅ Generic types require **type arguments** when instantiated
- ✅ Cannot use type parameters with **methods** (only functions and types)

---

## 🎯 Key Learning Objectives

After studying this code, you should understand:

1. **Functions**: How to write flexible, reusable functions using closures, higher-order functions, and variadic parameters
2. **Maps**: Efficient data storage and retrieval using hash tables, reference semantics, and advanced patterns
3. **Interfaces**: Polymorphic programming without inheritance, type assertions, and interface composition
4. **Generics**: Writing type-safe, reusable code with constraints and generic data structures

## 🔍 Code Organization

Each file follows this structure:
1. **Concept Sections**: Grouped by related functionality
2. **Detailed Comments**: Every important concept is explained
3. **Working Examples**: All code is tested and functional
4. **Demonstrate Function**: Each file has a `DemonstrateXXX()` function showcasing all features

## 📝 Best Practices Demonstrated

- ✅ **Error Handling**: Multiple return values for errors
- ✅ **Type Safety**: Using generics and interfaces appropriately
- ✅ **Code Organization**: Clear sections with comprehensive comments
- ✅ **Idiomatic Go**: Following Go conventions and patterns
- ✅ **Documentation**: Every function and type is documented

## 🛠️ Dependencies

This project uses:
- **golang.org/x/exp/constraints**: For the `Ordered` constraint in generics

Install dependencies:
```bash
go mod download
```

## 📖 Additional Resources

- [Go Documentation - Functions](https://go.dev/doc/effective_go#functions)
- [Go Documentation - Maps](https://go.dev/blog/maps)
- [Go Documentation - Interfaces](https://go.dev/doc/effective_go#interfaces)
- [Go Documentation - Generics](https://go.dev/doc/tutorial/generics)

## 🎓 Learning Path

**Recommended Order:**
1. Start with **functions.go** - Foundation for everything else
2. Move to **maps.go** - Essential data structure
3. Study **interfaces.go** - Key to polymorphism in Go
4. Finish with **generics.go** - Modern type-safe programming

---

**Happy Learning! 🚀**
