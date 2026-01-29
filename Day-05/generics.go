package main

import (
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

// =============================================================================
// TYPE PARAMETERS BASICS
// =============================================================================

// Generic function with type parameter
// T is a type parameter that can be any type
func Print[T any](value T) {
	fmt.Printf("Value: %v, Type: %T\n", value, value)
}

// Generic function with multiple type parameters
func MakePair[T, U any](first T, second U) (T, U) {
	return first, second
}

// Generic swap function
func Swap[T any](a, b T) (T, T) {
	return b, a
}

// =============================================================================
// CONSTRAINTS - COMPARABLE
// =============================================================================

// Function using comparable constraint
// comparable allows ==, != operations
func Contains[T comparable](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// Find index of element in slice
func IndexOf[T comparable](slice []T, value T) int {
	for i, item := range slice {
		if item == value {
			return i
		}
	}
	return -1
}

// Remove duplicates from slice
func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	result := []T{}

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// =============================================================================
// CONSTRAINTS - ORDERED
// =============================================================================

// Min function using ordered constraint from constraints package
// Ordered allows <, <=, >, >= operations
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max function
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Find minimum in slice
func MinSlice[T constraints.Ordered](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}

	min := slice[0]
	for _, item := range slice[1:] {
		if item < min {
			min = item
		}
	}
	return min
}

// Find maximum in slice
func MaxSlice[T constraints.Ordered](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}

	max := slice[0]
	for _, item := range slice[1:] {
		if item > max {
			max = item
		}
	}
	return max
}

// Clamp value between min and max
func Clamp[T constraints.Ordered](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// =============================================================================
// CUSTOM CONSTRAINTS
// =============================================================================

// Custom constraint using interface
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Sum function using custom Number constraint
func Sum[T Number](numbers []T) T {
	var sum T
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// Average function
func Average[T Number](numbers []T) float64 {
	if len(numbers) == 0 {
		return 0
	}

	var sum T
	for _, num := range numbers {
		sum += num
	}
	return float64(sum) / float64(len(numbers))
}

// Custom constraint for types with String() method
type Stringer interface {
	String() string
}

// Function accepting any type with String() method
func PrintStrings[T Stringer](items []T) {
	for _, item := range items {
		fmt.Println(item.String())
	}
}

// =============================================================================
// GENERIC STACK
// =============================================================================

// Stack - generic LIFO data structure
type Stack[T any] struct {
	items []T
}

// NewStack creates a new stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: []T{},
	}
}

// Push adds an item to the stack
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top item
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}

	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index]
	return item, true
}

// Peek returns the top item without removing it
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// IsEmpty checks if stack is empty
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of items
func (s *Stack[T]) Size() int {
	return len(s.items)
}

// =============================================================================
// GENERIC QUEUE
// =============================================================================

// Queue - generic FIFO data structure
type Queue[T any] struct {
	items []T
}

// NewQueue creates a new queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: []T{},
	}
}

// Enqueue adds an item to the queue
func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the front item
func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Front returns the front item without removing it
func (q *Queue[T]) Front() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	return q.items[0], true
}

// IsEmpty checks if queue is empty
func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// Size returns the number of items
func (q *Queue[T]) Size() int {
	return len(q.items)
}

// =============================================================================
// GENERIC LINKED LIST
// =============================================================================

// Node in a linked list
type Node[T any] struct {
	Value T
	Next  *Node[T]
}

// LinkedList - generic singly linked list
type LinkedList[T any] struct {
	Head *Node[T]
	size int
}

// NewLinkedList creates a new linked list
func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

// Append adds an item to the end
func (l *LinkedList[T]) Append(value T) {
	newNode := &Node[T]{Value: value}

	if l.Head == nil {
		l.Head = newNode
		l.size = 1
		return
	}

	current := l.Head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newNode
	l.size++
}

// Prepend adds an item to the beginning
func (l *LinkedList[T]) Prepend(value T) {
	newNode := &Node[T]{Value: value, Next: l.Head}
	l.Head = newNode
	l.size++
}

// Delete removes first occurrence of value
func (l *LinkedList[T]) Delete(value T, equals func(T, T) bool) bool {
	if l.Head == nil {
		return false
	}

	// Check if head needs to be deleted
	if equals(l.Head.Value, value) {
		l.Head = l.Head.Next
		l.size--
		return true
	}

	current := l.Head
	for current.Next != nil {
		if equals(current.Next.Value, value) {
			current.Next = current.Next.Next
			l.size--
			return true
		}
		current = current.Next
	}

	return false
}

// ToSlice converts linked list to slice
func (l *LinkedList[T]) ToSlice() []T {
	result := make([]T, 0, l.size)
	current := l.Head
	for current != nil {
		result = append(result, current.Value)
		current = current.Next
	}
	return result
}

// Size returns the number of items
func (l *LinkedList[T]) Size() int {
	return l.size
}

// =============================================================================
// GENERIC OPTIONAL TYPE
// =============================================================================

// Optional - represents a value that may or may not be present
type Optional[T any] struct {
	value   T
	present bool
}

// Some creates an Optional with a value
func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, present: true}
}

// None creates an empty Optional
func None[T any]() Optional[T] {
	return Optional[T]{present: false}
}

// IsPresent checks if value exists
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// Get returns the value (use after checking IsPresent)
func (o Optional[T]) Get() T {
	return o.value
}

// GetOrElse returns value or default if not present
func (o Optional[T]) GetOrElse(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

// Map transforms the value if present
func (o Optional[T]) Map(fn func(T) T) Optional[T] {
	if o.present {
		return Some(fn(o.value))
	}
	return None[T]()
}

// =============================================================================
// GENERIC RESULT TYPE (for error handling)
// =============================================================================

// Result - represents either success (Ok) or failure (Err)
type Result[T any] struct {
	value T
	err   error
	isOk  bool
}

// Ok creates a successful Result
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, isOk: true}
}

// Err creates a failed Result
func Err[T any](err error) Result[T] {
	return Result[T]{err: err, isOk: false}
}

// IsOk checks if result is successful
func (r Result[T]) IsOk() bool {
	return r.isOk
}

// IsErr checks if result is error
func (r Result[T]) IsErr() bool {
	return !r.isOk
}

// Unwrap returns value (panics if error)
func (r Result[T]) Unwrap() T {
	if !r.isOk {
		panic(r.err)
	}
	return r.value
}

// UnwrapOr returns value or default if error
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.isOk {
		return r.value
	}
	return defaultValue
}

// Error returns the error
func (r Result[T]) Error() error {
	return r.err
}

// =============================================================================
// GENERIC MAP/FILTER/REDUCE
// =============================================================================

// Map applies function to each element and returns new slice
func Map[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = fn(item)
	}
	return result
}

// Filter returns slice with elements that satisfy predicate
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := []T{}
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Reduce reduces slice to single value
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	result := initial
	for _, item := range slice {
		result = reducer(result, item)
	}
	return result
}

// FlatMap applies function and flattens result
func FlatMap[T, U any](slice []T, fn func(T) []U) []U {
	result := []U{}
	for _, item := range slice {
		result = append(result, fn(item)...)
	}
	return result
}

// GroupBy groups elements by key function
func GroupBy[T any, K comparable](slice []T, keyFn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range slice {
		key := keyFn(item)
		result[key] = append(result[key], item)
	}
	return result
}

// Partition splits slice into two based on predicate
func Partition[T any](slice []T, predicate func(T) bool) ([]T, []T) {
	trueSlice := []T{}
	falseSlice := []T{}

	for _, item := range slice {
		if predicate(item) {
			trueSlice = append(trueSlice, item)
		} else {
			falseSlice = append(falseSlice, item)
		}
	}

	return trueSlice, falseSlice
}

// All checks if all elements satisfy predicate
func All[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Any checks if any element satisfies predicate
func Any[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

// =============================================================================
// GENERIC PAIR/TUPLE
// =============================================================================

// Pair holds two values of potentially different types
type Pair[T, U any] struct {
	First  T
	Second U
}

// NewPair creates a new pair
func NewPair[T, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{First: first, Second: second}
}

// Zip combines two slices into slice of pairs
func Zip[T, U any](slice1 []T, slice2 []U) []Pair[T, U] {
	length := len(slice1)
	if len(slice2) < length {
		length = len(slice2)
	}

	result := make([]Pair[T, U], length)
	for i := 0; i < length; i++ {
		result[i] = NewPair(slice1[i], slice2[i])
	}
	return result
}

// =============================================================================
// DEMONSTRATION FUNCTION
// =============================================================================

// DemonstrateGenerics showcases all generic concepts
func DemonstrateGenerics() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("GENERICS DEMONSTRATION")
	fmt.Println(strings.Repeat("=", 80))

	// 1. Basic Type Parameters
	fmt.Println("\n1. BASIC TYPE PARAMETERS:")
	Print(42)
	Print("Hello, Generics!")
	Print(3.14)
	Print(true)

	first, second := MakePair(100, "hundred")
	fmt.Printf("MakePair: (%v, %v)\n", first, second)

	a, b := Swap("first", "second")
	fmt.Printf("Swap: %v, %v\n", a, b)

	// 2. Comparable Constraint
	fmt.Println("\n2. COMPARABLE CONSTRAINT:")

	numbers := []int{1, 2, 3, 4, 5}
	fmt.Printf("Contains(numbers, 3): %t\n", Contains(numbers, 3))
	fmt.Printf("Contains(numbers, 10): %t\n", Contains(numbers, 10))

	words := []string{"apple", "banana", "cherry"}
	fmt.Printf("IndexOf(words, 'banana'): %d\n", IndexOf(words, "banana"))

	duplicates := []int{1, 2, 2, 3, 3, 3, 4}
	unique := RemoveDuplicates(duplicates)
	fmt.Printf("RemoveDuplicates(%v) = %v\n", duplicates, unique)

	// 3. Ordered Constraint
	fmt.Println("\n3. ORDERED CONSTRAINT:")

	fmt.Printf("Min(10, 20) = %d\n", Min(10, 20))
	fmt.Printf("Max(10, 20) = %d\n", Max(10, 20))
	fmt.Printf("Min(3.14, 2.71) = %.2f\n", Min(3.14, 2.71))

	nums := []int{5, 2, 8, 1, 9, 3}
	fmt.Printf("MinSlice(%v) = %d\n", nums, MinSlice(nums))
	fmt.Printf("MaxSlice(%v) = %d\n", nums, MaxSlice(nums))

	fmt.Printf("Clamp(15, 0, 10) = %d\n", Clamp(15, 0, 10))
	fmt.Printf("Clamp(-5, 0, 10) = %d\n", Clamp(-5, 0, 10))
	fmt.Printf("Clamp(5, 0, 10) = %d\n", Clamp(5, 0, 10))

	// 4. Custom Constraints
	fmt.Println("\n4. CUSTOM CONSTRAINTS (Number):")

	intSlice := []int{1, 2, 3, 4, 5}
	fmt.Printf("Sum(int): %d\n", Sum(intSlice))
	fmt.Printf("Average(int): %.2f\n", Average(intSlice))

	floatSlice := []float64{1.5, 2.5, 3.5}
	fmt.Printf("Sum(float64): %.2f\n", Sum(floatSlice))
	fmt.Printf("Average(float64): %.2f\n", Average(floatSlice))

	// 5. Generic Stack
	fmt.Println("\n5. GENERIC STACK:")

	intStack := NewStack[int]()
	intStack.Push(10)
	intStack.Push(20)
	intStack.Push(30)

	fmt.Printf("Stack size: %d\n", intStack.Size())

	if val, ok := intStack.Peek(); ok {
		fmt.Printf("Peek: %d\n", val)
	}

	for !intStack.IsEmpty() {
		if val, ok := intStack.Pop(); ok {
			fmt.Printf("Pop: %d\n", val)
		}
	}

	// String stack
	stringStack := NewStack[string]()
	stringStack.Push("first")
	stringStack.Push("second")
	stringStack.Push("third")

	fmt.Printf("\nString stack size: %d\n", stringStack.Size())
	for !stringStack.IsEmpty() {
		if val, ok := stringStack.Pop(); ok {
			fmt.Printf("Pop: %s\n", val)
		}
	}

	// 6. Generic Queue
	fmt.Println("\n6. GENERIC QUEUE:")

	queue := NewQueue[string]()
	queue.Enqueue("first")
	queue.Enqueue("second")
	queue.Enqueue("third")

	fmt.Printf("Queue size: %d\n", queue.Size())

	for !queue.IsEmpty() {
		if val, ok := queue.Dequeue(); ok {
			fmt.Printf("Dequeue: %s\n", val)
		}
	}

	// 7. Generic LinkedList
	fmt.Println("\n7. GENERIC LINKED LIST:")

	list := NewLinkedList[int]()
	list.Append(1)
	list.Append(2)
	list.Append(3)
	list.Prepend(0)

	fmt.Printf("LinkedList: %v\n", list.ToSlice())
	fmt.Printf("Size: %d\n", list.Size())

	equals := func(a, b int) bool { return a == b }
	list.Delete(2, equals)
	fmt.Printf("After deleting 2: %v\n", list.ToSlice())

	// 8. Optional Type
	fmt.Println("\n8. OPTIONAL TYPE:")

	opt1 := Some(42)
	opt2 := None[int]()

	fmt.Printf("opt1.IsPresent(): %t, Value: %d\n", opt1.IsPresent(), opt1.Get())
	fmt.Printf("opt2.IsPresent(): %t\n", opt2.IsPresent())

	fmt.Printf("opt1.GetOrElse(0): %d\n", opt1.GetOrElse(0))
	fmt.Printf("opt2.GetOrElse(100): %d\n", opt2.GetOrElse(100))

	mapped := opt1.Map(func(x int) int { return x * 2 })
	fmt.Printf("opt1.Map(×2): %d\n", mapped.Get())

	// 9. Result Type
	fmt.Println("\n9. RESULT TYPE:")

	res1 := Ok(100)
	res2 := Err[int](fmt.Errorf("something went wrong"))

	fmt.Printf("res1.IsOk(): %t, Value: %d\n", res1.IsOk(), res1.Unwrap())
	fmt.Printf("res2.IsErr(): %t, Error: %v\n", res2.IsErr(), res2.Error())

	fmt.Printf("res1.UnwrapOr(0): %d\n", res1.UnwrapOr(0))
	fmt.Printf("res2.UnwrapOr(999): %d\n", res2.UnwrapOr(999))

	// 10. Map/Filter/Reduce
	fmt.Println("\n10. MAP/FILTER/REDUCE:")

	input := []int{1, 2, 3, 4, 5}

	// Map
	squared := Map(input, func(x int) int { return x * x })
	fmt.Printf("Map (square): %v -> %v\n", input, squared)

	// Filter
	evens := Filter(input, func(x int) bool { return x%2 == 0 })
	fmt.Printf("Filter (evens): %v -> %v\n", input, evens)

	// Reduce
	sum := Reduce(input, 0, func(acc, x int) int { return acc + x })
	fmt.Printf("Reduce (sum): %v -> %d\n", input, sum)

	// FlatMap
	doubled := FlatMap(input, func(x int) []int { return []int{x, x} })
	fmt.Printf("FlatMap (duplicate): %v -> %v\n", input, doubled)

	// GroupBy
	grouped := GroupBy(input, func(x int) string {
		if x%2 == 0 {
			return "even"
		}
		return "odd"
	})
	fmt.Printf("GroupBy (even/odd): %v\n", grouped)

	// Partition
	evenPart, oddPart := Partition(input, func(x int) bool { return x%2 == 0 })
	fmt.Printf("Partition: evens=%v, odds=%v\n", evenPart, oddPart)

	// All/Any
	allPositive := All(input, func(x int) bool { return x > 0 })
	anyNegative := Any(input, func(x int) bool { return x < 0 })
	fmt.Printf("All positive: %t, Any negative: %t\n", allPositive, anyNegative)

	// 11. Pair and Zip
	fmt.Println("\n11. PAIR AND ZIP:")

	pair := NewPair(42, "answer")
	fmt.Printf("Pair: (%v, %v)\n", pair.First, pair.Second)

	names := []string{"Alice", "Bob", "Charlie"}
	ages := []int{30, 25, 35}
	zipped := Zip(names, ages)

	fmt.Println("Zipped:")
	for _, p := range zipped {
		fmt.Printf("  %s: %d\n", p.First, p.Second)
	}

	// 12. Chaining Operations
	fmt.Println("\n12. CHAINING OPERATIONS:")

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Filter evens, map to squares, sum
	result := Reduce(
		Map(
			Filter(data, func(x int) bool { return x%2 == 0 }),
			func(x int) int { return x * x },
		),
		0,
		func(acc, x int) int { return acc + x },
	)

	fmt.Printf("Sum of squares of evens in %v: %d\n", data, result)

	fmt.Println("\n" + strings.Repeat("=", 80))
}
