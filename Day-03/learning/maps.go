package learning

import "fmt"

// ==========================================
// MAPS IN GO
// ==========================================
// Maps are key-value data structures (like dictionaries/hash tables)
// Syntax: map[KeyType]ValueType

// DemonstrateMapBasics shows basic map operations
func DemonstrateMapBasics() {
	fmt.Println("\n========== MAP BASICS ==========")

	// 1. Declaration and Initialization
	// Method 1: Using make function (recommended)
	var scores map[string]int = make(map[string]int)
	scores["Alice"] = 95
	scores["Bob"] = 87
	scores["Charlie"] = 92

	// Method 2: Map literal (initialize with values)
	grades := map[string]string{
		"Alice":   "A",
		"Bob":     "B",
		"Charlie": "A-",
	}

	// Method 3: Declare without make (will be nil)
	var nilMap map[string]int // nil map, can't add to it
	fmt.Printf("Nil map: %v, is nil: %v\n", nilMap, nilMap == nil)

	fmt.Printf("Scores: %v\n", scores)
	fmt.Printf("Grades: %v\n", grades)

	// 2. Accessing elements
	aliceScore := scores["Alice"]
	fmt.Printf("Alice's score: %d\n", aliceScore)

	// 3. Check if key exists (very important!)
	score, exists := scores["David"]
	if exists {
		fmt.Printf("David's score: %d\n", score)
	} else {
		fmt.Println("David not found in map")
		fmt.Printf("Default value returned: %d\n", score) // 0 for int
	}

	// 4. Adding/Updating elements
	scores["David"] = 88 // Add new
	scores["Alice"] = 98 // Update existing
	fmt.Printf("Updated scores: %v\n", scores)

	// 5. Deleting elements
	delete(scores, "Bob")
	fmt.Printf("After deleting Bob: %v\n", scores)

	// 6. Length of map
	fmt.Printf("Number of students: %d\n", len(scores))
}

// DemonstrateMapIteration shows different ways to iterate over maps
func DemonstrateMapIteration() {
	fmt.Println("\n========== MAP ITERATION ==========")

	countries := map[string]string{
		"US": "United States",
		"UK": "United Kingdom",
		"IN": "India",
		"JP": "Japan",
		"FR": "France",
	}

	// 1. Iterate over key-value pairs
	fmt.Println("Countries:")
	for code, name := range countries {
		fmt.Printf("  %s => %s\n", code, name)
	}

	// 2. Iterate over keys only
	fmt.Println("\nCountry Codes:")
	for code := range countries {
		fmt.Printf("  %s\n", code)
	}

	// 3. Iterate over values only (use _ for key)
	fmt.Println("\nCountry Names:")
	for _, name := range countries {
		fmt.Printf("  %s\n", name)
	}

	// NOTE: Map iteration order is NOT guaranteed!
	// Running the same code twice may give different order
}

// DemonstrateComplexMaps shows maps with complex types
func DemonstrateComplexMaps() {
	fmt.Println("\n========== COMPLEX MAPS ==========")

	// 1. Map with slice as value
	studentCourses := map[string][]string{
		"Alice":   {"Math", "Physics", "Chemistry"},
		"Bob":     {"Biology", "English"},
		"Charlie": {"Math", "Computer Science", "Physics"},
	}

	fmt.Println("Student Courses:")
	for student, courses := range studentCourses {
		fmt.Printf("%s: %v\n", student, courses)
	}

	// 2. Map with struct as value (we'll define simple inline struct)
	type Student struct {
		Age   int
		Grade string
	}

	students := map[string]Student{
		"Alice": {Age: 20, Grade: "A"},
		"Bob":   {Age: 21, Grade: "B"},
	}

	fmt.Println("\nStudent Details:")
	for name, details := range students {
		fmt.Printf("%s: Age=%d, Grade=%s\n", name, details.Age, details.Grade)
	}

	// 3. Nested maps (map of maps)
	// Storing test scores by subject for each student
	testScores := map[string]map[string]int{
		"Alice": {
			"Math":    95,
			"Physics": 92,
			"English": 88,
		},
		"Bob": {
			"Math":    78,
			"Physics": 82,
			"English": 90,
		},
	}

	fmt.Println("\nTest Scores (Nested Map):")
	for student, subjects := range testScores {
		fmt.Printf("%s:\n", student)
		for subject, score := range subjects {
			fmt.Printf("  %s: %d\n", subject, score)
		}
	}

	// Adding to nested map
	testScores["Charlie"] = make(map[string]int)
	testScores["Charlie"]["Math"] = 85
	testScores["Charlie"]["Physics"] = 88
}

// DemonstrateMapAsReference shows that maps are reference types
func DemonstrateMapAsReference() {
	fmt.Println("\n========== MAPS ARE REFERENCE TYPES ==========")

	// Maps are reference types (like slices)
	// When you assign a map to another variable, both refer to the same data
	original := map[string]int{
		"a": 1,
		"b": 2,
	}

	copied := original // This doesn't create a new map, just another reference
	copied["c"] = 3    // This affects the original map too!

	fmt.Printf("Original map: %v\n", original) // Will include "c"
	fmt.Printf("Copied map: %v\n", copied)

	// To create a true copy, you need to do it manually
	realCopy := make(map[string]int)
	for k, v := range original {
		realCopy[k] = v
	}
	realCopy["d"] = 4

	fmt.Printf("Original after real copy modification: %v\n", original) // No "d"
	fmt.Printf("Real copy: %v\n", realCopy)                             // Has "d"
}

// DemonstrateMapPracticalExample shows a practical use case
func DemonstrateMapPracticalExample() {
	fmt.Println("\n========== PRACTICAL EXAMPLE: Word Counter ==========")

	wordCount := make(map[string]int)

	// Simple word counting (splitting by space)
	words := []string{"the", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog", "the", "fox", "is", "quick"}

	for _, word := range words {
		wordCount[word]++ // If key doesn't exist, default value (0) is used, then incremented
	}

	fmt.Println("Word frequencies:")
	for word, count := range wordCount {
		fmt.Printf("  '%s': %d times\n", word, count)
	}
}

// RunAllMapDemos runs all map demonstrations
func RunAllMapDemos() {
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║     LEARNING MAPS IN GO               ║")
	fmt.Println("╚════════════════════════════════════════╝")

	DemonstrateMapBasics()
	DemonstrateMapIteration()
	DemonstrateComplexMaps()
	DemonstrateMapAsReference()
	DemonstrateMapPracticalExample()
}

