package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("GO ADVANCED CONCEPTS - DAY 05")
	fmt.Println(strings.Repeat("=", 60))

	// ============================================
	// Type Conversion Examples
	// ============================================
	fmt.Println("\n--- Type Conversion ---")
	i := 42
	f := float64(i)
	f = f + 0.99
	u := uint(f)
	fmt.Printf("Integer: %d, Float: %.2f, Unsigned: %d\n", i, f, u)

	// ============================================
	// Complex Numbers in Go
	// ============================================
	fmt.Println("\n--- Complex Numbers ---")
	c1 := complex(2, 3)
	c2 := 4 + 5i             // complex initializer syntax a + ib
	c3 := c1 + c2            // addition just like other variables
	fmt.Println("Add: ", c3) // prints "Add: (6+8i)"
	re := real(c3)           // get real part
	im := imag(c3)           // get imaginary part
	fmt.Printf("Real part: %.0f, Imaginary part: %.0f\n", re, im)

	// ============================================
	// Call all demonstration functions
	// ============================================

	// 1. Functions - Comprehensive coverage of all function concepts
	DemonstrateFunctions()

	// 2. Maps - All map operations and use cases
	DemonstrateMaps()

	// 3. Interfaces - Interface concepts and polymorphism
	DemonstrateInterfaces()

	// 4. Generics - Type parameters and generic programming
	DemonstrateGenerics()

	// ============================================
	// Summary
	// ============================================
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("ALL DEMONSTRATIONS COMPLETED!")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\nTopics Covered:")
	fmt.Println("  ✓ Functions (Basics, Variadic, Multiple Returns, Anonymous, Closures, Named Returns, Call by Value)")
	fmt.Println("  ✓ Maps (Creation, Operations, Iteration, Complex Types, Reference Semantics)")
	fmt.Println("  ✓ Interfaces (Definition, Implementation, Type Assertion, Polymorphism, Composition)")
	fmt.Println("  ✓ Generics (Type Parameters, Constraints, Generic Types, Advanced Patterns)")
	fmt.Println()
}
