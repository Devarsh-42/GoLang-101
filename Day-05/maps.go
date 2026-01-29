package main

import (
	"fmt"
	"sort"
	"strings"
)

// =============================================================================
// MAP CREATION METHODS
// =============================================================================

// Method 1: Using make() - most common and preferred
func createMapWithMake() map[string]int {
	// Create an empty map
	ages := make(map[string]int)
	ages["Alice"] = 30
	ages["Bob"] = 25
	return ages
}

// Method 2: Map literal - initialize with values
func createMapWithLiteral() map[string]int {
	return map[string]int{
		"Alice": 30,
		"Bob":   25,
		"Carol": 35,
	}
}

// Method 3: Empty map literal
func createEmptyMapLiteral() map[string]int {
	return map[string]int{}
}

// Method 4: Declare without initialization (nil map)
// WARNING: Cannot add to nil map - will panic
func demonstrateNilMap() {
	var m map[string]int // nil map
	fmt.Printf("Nil map: %v, len: %d, is nil: %v\n", m, len(m), m == nil)

	// Reading from nil map is safe - returns zero value
	val := m["key"]
	fmt.Printf("Reading from nil map: %d\n", val)

	// Writing to nil map will PANIC - uncomment to see
	// m["key"] = 10 // panic: assignment to entry in nil map
}

// =============================================================================
// BASIC MAP OPERATIONS
// =============================================================================

// Adding and updating elements
func addUpdateElements() {
	m := make(map[string]int)

	// Adding elements
	m["one"] = 1
	m["two"] = 2
	m["three"] = 3

	// Updating elements (same syntax as adding)
	m["one"] = 100

	fmt.Printf("Map after add/update: %v\n", m)
}

// Deleting elements
func deleteElements() {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Printf("Before delete: %v\n", m)

	// Delete key "b"
	delete(m, "b")
	fmt.Printf("After delete: %v\n", m)

	// Deleting non-existent key is safe (no-op)
	delete(m, "nonexistent")
	fmt.Printf("After deleting non-existent key: %v\n", m)
}

// Checking if key exists
func checkKeyExists() {
	m := map[string]int{"Alice": 30, "Bob": 25}

	// Two-value assignment to check if key exists
	value, exists := m["Alice"]
	if exists {
		fmt.Printf("Alice exists with value: %d\n", value)
	}

	value, exists = m["Charlie"]
	if !exists {
		fmt.Printf("Charlie does not exist (value: %d, exists: %v)\n", value, exists)
	}

	// Common pattern: checking and using in one statement
	if age, ok := m["Bob"]; ok {
		fmt.Printf("Bob's age: %d\n", age)
	}
}

// Getting map length
func getMapLength() {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Printf("Map: %v, Length: %d\n", m, len(m))
}

// =============================================================================
// MAP ITERATION
// =============================================================================

// Basic iteration over map
func basicIteration() {
	scores := map[string]int{
		"Alice":   95,
		"Bob":     87,
		"Charlie": 92,
		"David":   88,
	}

	fmt.Println("\nIterating over map (key-value pairs):")
	for name, score := range scores {
		fmt.Printf("%s: %d\n", name, score)
	}

	// Iterate over keys only
	fmt.Println("\nIterating over keys only:")
	for name := range scores {
		fmt.Printf("%s\n", name)
	}

	// Iterate over values only (using blank identifier)
	fmt.Println("\nIterating over values only:")
	total := 0
	for _, score := range scores {
		total += score
	}
	fmt.Printf("Total score: %d\n", total)
}

// Sorted iteration (maps are unordered)
func sortedIteration() {
	m := map[string]int{
		"banana": 3,
		"apple":  5,
		"cherry": 2,
		"date":   4,
	}

	// Extract keys and sort them
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fmt.Println("\nSorted iteration by keys:")
	for _, key := range keys {
		fmt.Printf("%s: %d\n", key, m[key])
	}
}

// =============================================================================
// MAPS WITH COMPLEX TYPES
// =============================================================================

// Map with struct values
type Person struct {
	Name  string
	Age   int
	Email string
}

func mapWithStructValues() map[int]Person {
	people := map[int]Person{
		1: {Name: "Alice", Age: 30, Email: "alice@example.com"},
		2: {Name: "Bob", Age: 25, Email: "bob@example.com"},
		3: {Name: "Carol", Age: 35, Email: "carol@example.com"},
	}
	return people
}

// Map with slice values
func mapWithSliceValues() {
	// Map of string to slice of ints
	grades := map[string][]int{
		"Alice":   {95, 87, 92},
		"Bob":     {88, 76, 91},
		"Charlie": {92, 89, 95},
	}

	fmt.Println("\nMap with slice values:")
	for name, scores := range grades {
		avg := 0.0
		for _, score := range scores {
			avg += float64(score)
		}
		avg /= float64(len(scores))
		fmt.Printf("%s: %v (avg: %.2f)\n", name, scores, avg)
	}

	// Adding to a slice value
	grades["Alice"] = append(grades["Alice"], 88)
	fmt.Printf("After adding score to Alice: %v\n", grades["Alice"])
}

// Map with map values (nested maps)
func mapWithMapValues() {
	// Map of maps - representing a 2D grid or hierarchical data
	config := map[string]map[string]string{
		"database": {
			"host":     "localhost",
			"port":     "5432",
			"username": "admin",
		},
		"cache": {
			"host": "localhost",
			"port": "6379",
			"ttl":  "3600",
		},
	}

	fmt.Println("\nNested maps:")
	for service, settings := range config {
		fmt.Printf("%s:\n", service)
		for key, value := range settings {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	// Accessing nested values
	dbHost := config["database"]["host"]
	fmt.Printf("Database host: %s\n", dbHost)
}

// Map with struct keys (must be comparable)
type Coordinate struct {
	X, Y int
}

func mapWithStructKeys() {
	// Structs can be map keys if all fields are comparable
	grid := map[Coordinate]string{
		{0, 0}: "origin",
		{1, 0}: "right",
		{0, 1}: "up",
		{1, 1}: "diagonal",
	}

	fmt.Println("\nMap with struct keys:")
	for coord, label := range grid {
		fmt.Printf("(%d, %d): %s\n", coord.X, coord.Y, label)
	}

	// Checking if coordinate exists
	coord := Coordinate{1, 0}
	if label, ok := grid[coord]; ok {
		fmt.Printf("Coordinate %v has label: %s\n", coord, label)
	}
}

// =============================================================================
// MAP REFERENCE SEMANTICS
// =============================================================================

// Maps are reference types - they refer to underlying data structure
func demonstrateReferenceSemantics() {
	// Original map
	original := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Printf("Original: %v\n", original)

	// Assignment creates a reference, not a copy
	reference := original
	reference["a"] = 100

	fmt.Printf("After modifying reference:\n")
	fmt.Printf("  Original: %v\n", original) // Also changed!
	fmt.Printf("  Reference: %v\n", reference)

	// To create a copy, you must manually copy elements
	copy := make(map[string]int)
	for k, v := range original {
		copy[k] = v
	}
	copy["a"] = 999

	fmt.Printf("After modifying copy:\n")
	fmt.Printf("  Original: %v\n", original) // Unchanged
	fmt.Printf("  Copy: %v\n", copy)
}

// Maps passed to functions are passed by reference
func modifyMap(m map[string]int) {
	m["modified"] = 100
}

func demonstrateMapAsParameter() {
	m := map[string]int{"a": 1, "b": 2}
	fmt.Printf("Before function call: %v\n", m)

	modifyMap(m)

	fmt.Printf("After function call: %v\n", m) // Map is modified
}

// =============================================================================
// ADVANCED MAP PATTERNS
// =============================================================================

// Using map as a set (set of unique values)
func mapAsSet() {
	// Use map[type]bool or map[type]struct{}
	// struct{} uses zero memory
	set := make(map[string]struct{})

	// Adding elements
	set["apple"] = struct{}{}
	set["banana"] = struct{}{}
	set["cherry"] = struct{}{}

	// Checking membership
	if _, exists := set["apple"]; exists {
		fmt.Println("apple is in the set")
	}

	// Adding duplicate has no effect
	set["apple"] = struct{}{}

	// Iterating over set
	fmt.Print("Set elements: ")
	for item := range set {
		fmt.Printf("%s ", item)
	}
	fmt.Println()
}

// Counting occurrences using map
func countOccurrences(words []string) map[string]int {
	counts := make(map[string]int)
	for _, word := range words {
		counts[word]++ // If key doesn't exist, it's initialized to 0
	}
	return counts
}

// Grouping items using map
func groupByFirstLetter(words []string) map[string][]string {
	groups := make(map[string][]string)

	for _, word := range words {
		if len(word) == 0 {
			continue
		}
		firstLetter := string(word[0])
		groups[firstLetter] = append(groups[firstLetter], word)
	}

	return groups
}

// Inverting a map (swap keys and values)
func invertMap(m map[string]int) map[int]string {
	inverted := make(map[int]string)
	for key, value := range m {
		// Note: if values are not unique, later entries will overwrite earlier ones
		inverted[value] = key
	}
	return inverted
}

// Merging maps
func mergeMaps(m1, m2 map[string]int) map[string]int {
	result := make(map[string]int)

	// Copy from first map
	for k, v := range m1 {
		result[k] = v
	}

	// Copy from second map (overwrites duplicates)
	for k, v := range m2 {
		result[k] = v
	}

	return result
}

// Map with default values
type DefaultMap struct {
	data         map[string]int
	defaultValue int
}

func NewDefaultMap(defaultValue int) *DefaultMap {
	return &DefaultMap{
		data:         make(map[string]int),
		defaultValue: defaultValue,
	}
}

func (dm *DefaultMap) Get(key string) int {
	if val, ok := dm.data[key]; ok {
		return val
	}
	return dm.defaultValue
}

func (dm *DefaultMap) Set(key string, value int) {
	dm.data[key] = value
}

// Synchronized map for concurrent access (simple example)
type SafeMap struct {
	data map[string]int
	// In real code, use sync.RWMutex for thread safety
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]int),
	}
}

// =============================================================================
// MAP PERFORMANCE AND CAPACITY
// =============================================================================

// Creating map with initial capacity hint
func mapWithCapacity() {
	// Providing size hint can improve performance for large maps
	m := make(map[string]int, 1000) // Hint: will store ~1000 elements

	for i := 0; i < 1000; i++ {
		m[fmt.Sprintf("key%d", i)] = i
	}

	fmt.Printf("Map with capacity hint - length: %d\n", len(m))
}

// =============================================================================
// DEMONSTRATION FUNCTION
// =============================================================================

// DemonstrateMaps showcases all map concepts
func DemonstrateMaps() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("MAPS DEMONSTRATION")
	fmt.Println(strings.Repeat("=", 80))

	// 1. Map Creation
	fmt.Println("\n1. MAP CREATION METHODS:")
	m1 := createMapWithMake()
	fmt.Printf("Created with make(): %v\n", m1)

	m2 := createMapWithLiteral()
	fmt.Printf("Created with literal: %v\n", m2)

	m3 := createEmptyMapLiteral()
	fmt.Printf("Empty map literal: %v\n", m3)

	demonstrateNilMap()

	// 2. Basic Operations
	fmt.Println("\n2. BASIC MAP OPERATIONS:")
	addUpdateElements()
	deleteElements()
	checkKeyExists()
	getMapLength()

	// 3. Map Iteration
	fmt.Println("\n3. MAP ITERATION:")
	basicIteration()
	sortedIteration()

	// 4. Maps with Complex Types
	fmt.Println("\n4. MAPS WITH COMPLEX TYPES:")

	people := mapWithStructValues()
	fmt.Println("Map with struct values:")
	for id, person := range people {
		fmt.Printf("  ID %d: %s (age %d)\n", id, person.Name, person.Age)
	}

	mapWithSliceValues()
	mapWithMapValues()
	mapWithStructKeys()

	// 5. Reference Semantics
	fmt.Println("\n5. MAP REFERENCE SEMANTICS:")
	demonstrateReferenceSemantics()
	demonstrateMapAsParameter()

	// 6. Advanced Patterns
	fmt.Println("\n6. ADVANCED MAP PATTERNS:")

	// Map as set
	fmt.Println("Map as Set:")
	mapAsSet()

	// Counting occurrences
	words := []string{"apple", "banana", "apple", "cherry", "banana", "apple"}
	counts := countOccurrences(words)
	fmt.Printf("\nWord counts: %v\n", counts)

	// Grouping
	wordList := []string{"apple", "apricot", "banana", "blueberry", "cherry", "cranberry"}
	groups := groupByFirstLetter(wordList)
	fmt.Println("\nGrouped by first letter:")
	for letter, words := range groups {
		fmt.Printf("  %s: %v\n", letter, words)
	}

	// Inverting map
	original := map[string]int{"one": 1, "two": 2, "three": 3}
	inverted := invertMap(original)
	fmt.Printf("\nOriginal: %v\n", original)
	fmt.Printf("Inverted: %v\n", inverted)

	// Merging maps
	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"b": 20, "c": 3}
	merged := mergeMaps(map1, map2)
	fmt.Printf("\nMerging %v and %v = %v\n", map1, map2, merged)

	// Default map
	fmt.Println("\nDefault Map:")
	dm := NewDefaultMap(42)
	dm.Set("exists", 100)
	fmt.Printf("Key 'exists': %d\n", dm.Get("exists"))
	fmt.Printf("Key 'missing': %d (default)\n", dm.Get("missing"))

	// 7. Performance
	fmt.Println("\n7. MAP CAPACITY:")
	mapWithCapacity()

	// 8. Common Patterns
	fmt.Println("\n8. COMMON MAP PATTERNS:")

	// Frequency counter
	text := "hello world hello"
	charCount := make(map[rune]int)
	for _, char := range text {
		if char != ' ' {
			charCount[char]++
		}
	}
	fmt.Printf("Character frequency: %v\n", charCount)

	// Cache/memoization
	cache := make(map[int]int)
	fibonacci := func(n int) int {
		if n <= 1 {
			return n
		}
		if val, ok := cache[n]; ok {
			return val
		}
		result := n // Simplified for demo
		cache[n] = result
		return result
	}
	fmt.Printf("Fibonacci cache example: %d\n", fibonacci(10))

	// Two-sum problem using map
	nums := []int{2, 7, 11, 15}
	target := 9
	seen := make(map[int]int)
	fmt.Printf("\nFinding two numbers that sum to %d in %v:\n", target, nums)
	for i, num := range nums {
		complement := target - num
		if j, ok := seen[complement]; ok {
			fmt.Printf("Found: %d + %d = %d (indices %d, %d)\n",
				nums[j], num, target, j, i)
			break
		}
		seen[num] = i
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
}
