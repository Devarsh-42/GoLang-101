package learning

import (
	"fmt"
	"math"
	"strings"
)

// ==========================================
// RECEIVER FUNCTIONS (METHODS) IN GO
// ==========================================
// Methods are functions with a receiver argument
// They allow you to attach behavior to types
// Syntax: func (receiver Type) MethodName(params) returnType

// ========== VALUE RECEIVERS ==========

// Rectangle demonstrates methods with value receivers
type Rectangle struct {
	Width  float64
	Height float64
}

// Area calculates the area of a rectangle (value receiver)
// Value receiver gets a copy of the struct
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter calculates the perimeter of a rectangle (value receiver)
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Info returns formatted information about the rectangle (value receiver)
func (r Rectangle) Info() string {
	return fmt.Sprintf("Rectangle: %.2fx%.2f (Area: %.2f, Perimeter: %.2f)",
		r.Width, r.Height, r.Area(), r.Perimeter())
}

// TryToScale attempts to scale the rectangle (value receiver - won't modify original!)
func (r Rectangle) TryToScale(factor float64) {
	r.Width *= factor
	r.Height *= factor
	fmt.Println("  [Inside TryToScale] Scaled dimensions:", r.Width, "x", r.Height)
}

// DemonstrateValueReceivers shows value receiver behavior
func DemonstrateValueReceivers() {
	fmt.Println("\n========== VALUE RECEIVERS ==========")

	rect := Rectangle{Width: 10, Height: 5}

	fmt.Printf("Original: %s\n", rect.Info())
	fmt.Printf("Area: %.2f\n", rect.Area())
	fmt.Printf("Perimeter: %.2f\n", rect.Perimeter())

	// Value receiver doesn't modify original
	fmt.Println("\nCalling TryToScale(2):")
	rect.TryToScale(2)
	fmt.Printf("After TryToScale: %s\n", rect.Info()) // Unchanged!

	fmt.Println("\n✓ Value receivers work on a COPY of the data")
	fmt.Println("✓ Use when: method doesn't need to modify the receiver")
	fmt.Println("✓ Use when: receiver is small (copying is cheap)")
}

// ========== POINTER RECEIVERS ==========

// Circle demonstrates methods with pointer receivers
type Circle struct {
	Radius float64
}

// Area calculates the area of a circle (pointer receiver)
// Pointer receiver gets access to the original struct
func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Circumference calculates the circumference of a circle (pointer receiver)
func (c *Circle) Circumference() float64 {
	return 2 * math.Pi * c.Radius
}

// Scale scales the circle radius (pointer receiver - WILL modify original)
func (c *Circle) Scale(factor float64) {
	c.Radius *= factor
	fmt.Printf("  [Inside Scale] New radius: %.2f\n", c.Radius)
}

// Reset resets the circle to a default radius (pointer receiver)
func (c *Circle) Reset() {
	c.Radius = 1.0
}

// DemonstratePointerReceivers shows pointer receiver behavior
func DemonstratePointerReceivers() {
	fmt.Println("\n========== POINTER RECEIVERS ==========")

	circle := Circle{Radius: 5.0}

	fmt.Printf("Original radius: %.2f\n", circle.Radius)
	fmt.Printf("Area: %.2f\n", circle.Area())
	fmt.Printf("Circumference: %.2f\n", circle.Circumference())

	// Pointer receiver DOES modify original
	fmt.Println("\nCalling Scale(2):")
	circle.Scale(2)
	fmt.Printf("After Scale: Radius = %.2f\n", circle.Radius) // Changed to 10!

	// Go automatically handles pointer/value conversions
	// Even though Scale expects *Circle, we can call it on Circle
	// Go automatically takes the address (&circle)

	fmt.Println("\n✓ Pointer receivers work on the ORIGINAL data")
	fmt.Println("✓ Use when: method needs to modify the receiver")
	fmt.Println("✓ Use when: receiver is large (avoid copying)")
	fmt.Println("✓ Use when: consistency (if one method uses pointer, all should)")
}

// ========== METHODS ON CUSTOM TYPES ==========

// Temperature is a custom type based on float64
type Temperature float64

// Celsius returns the temperature in Celsius
func (t Temperature) Celsius() float64 {
	return float64(t)
}

// Fahrenheit converts Celsius to Fahrenheit
func (t Temperature) Fahrenheit() float64 {
	return float64(t)*9/5 + 32
}

// Kelvin converts Celsius to Kelvin
func (t Temperature) Kelvin() float64 {
	return float64(t) + 273.15
}

// IsFreezing checks if temperature is at or below freezing
func (t Temperature) IsFreezing() bool {
	return t <= 0
}

// IsBoiling checks if temperature is at or above boiling point
func (t Temperature) IsBoiling() bool {
	return t >= 100
}

// Description returns a description of the temperature
func (t Temperature) Description() string {
	switch {
	case t < 0:
		return "Freezing"
	case t < 10:
		return "Cold"
	case t < 20:
		return "Cool"
	case t < 30:
		return "Warm"
	default:
		return "Hot"
	}
}

// DemonstrateMethodsOnCustomTypes shows methods on custom types
func DemonstrateMethodsOnCustomTypes() {
	fmt.Println("\n========== METHODS ON CUSTOM TYPES ==========")

	temp := Temperature(25.0) // 25°C

	fmt.Printf("Temperature: %.2f°C\n", temp.Celsius())
	fmt.Printf("In Fahrenheit: %.2f°F\n", temp.Fahrenheit())
	fmt.Printf("In Kelvin: %.2fK\n", temp.Kelvin())
	fmt.Printf("Description: %s\n", temp.Description())
	fmt.Printf("Is Freezing: %v\n", temp.IsFreezing())
	fmt.Printf("Is Boiling: %v\n", temp.IsBoiling())

	// Test with freezing temperature
	freezing := Temperature(-5.0)
	fmt.Printf("\nAt %.2f°C:\n", freezing.Celsius())
	fmt.Printf("Description: %s\n", freezing.Description())
	fmt.Printf("Is Freezing: %v\n", freezing.IsFreezing())
}

// ========== METHODS WITH PARAMETERS ==========

// BankAccount demonstrates methods with parameters
type BankAccount struct {
	AccountNumber string
	HolderName    string
	Balance       float64
}

// Deposit adds money to the account (pointer receiver to modify balance)
func (b *BankAccount) Deposit(amount float64) {
	if amount > 0 {
		b.Balance += amount
		fmt.Printf("✓ Deposited $%.2f. New balance: $%.2f\n", amount, b.Balance)
	} else {
		fmt.Println("✗ Deposit amount must be positive")
	}
}

// Withdraw removes money from the account (pointer receiver to modify balance)
func (b *BankAccount) Withdraw(amount float64) bool {
	if amount <= 0 {
		fmt.Println("✗ Withdrawal amount must be positive")
		return false
	}
	if amount > b.Balance {
		fmt.Printf("✗ Insufficient funds. Balance: $%.2f, Requested: $%.2f\n", b.Balance, amount)
		return false
	}
	b.Balance -= amount
	fmt.Printf("✓ Withdrew $%.2f. New balance: $%.2f\n", amount, b.Balance)
	return true
}

// Transfer transfers money to another account (pointer receiver)
func (b *BankAccount) Transfer(to *BankAccount, amount float64) bool {
	fmt.Printf("Transferring $%.2f from %s to %s\n", amount, b.HolderName, to.HolderName)
	if b.Withdraw(amount) {
		to.Deposit(amount)
		return true
	}
	return false
}

// GetBalance returns the current balance (value receiver - read-only)
func (b BankAccount) GetBalance() float64 {
	return b.Balance
}

// Info returns account information (value receiver - read-only)
func (b BankAccount) Info() string {
	return fmt.Sprintf("Account: %s | Holder: %s | Balance: $%.2f",
		b.AccountNumber, b.HolderName, b.Balance)
}

// DemonstrateMethodsWithParameters shows methods with parameters
func DemonstrateMethodsWithParameters() {
	fmt.Println("\n========== METHODS WITH PARAMETERS ==========")

	account1 := BankAccount{
		AccountNumber: "ACC001",
		HolderName:    "Alice",
		Balance:       1000.00,
	}

	account2 := BankAccount{
		AccountNumber: "ACC002",
		HolderName:    "Bob",
		Balance:       500.00,
	}

	fmt.Println("Initial State:")
	fmt.Println(account1.Info())
	fmt.Println(account2.Info())

	fmt.Println("\nOperations:")
	account1.Deposit(200)
	account1.Withdraw(150)
	account1.Withdraw(2000) // Should fail

	fmt.Println("\nTransfer:")
	account1.Transfer(&account2, 300)

	fmt.Println("\nFinal State:")
	fmt.Println(account1.Info())
	fmt.Println(account2.Info())
}

// ========== METHODS VS FUNCTIONS ==========

// CalculateRectangleArea is a regular function (not a method)
func CalculateRectangleArea(width, height float64) float64 {
	return width * height
}

// DemonstrateMethodsVsFunctions shows the difference
func DemonstrateMethodsVsFunctions() {
	fmt.Println("\n========== METHODS VS FUNCTIONS ==========")

	rect := Rectangle{Width: 10, Height: 5}

	// Using a method (called on a receiver)
	areaMethod := rect.Area()
	fmt.Printf("Using method: rect.Area() = %.2f\n", areaMethod)

	// Using a function (standalone)
	areaFunction := CalculateRectangleArea(rect.Width, rect.Height)
	fmt.Printf("Using function: CalculateRectangleArea(10, 5) = %.2f\n", areaFunction)

	fmt.Println("\nDifferences:")
	fmt.Println("• Methods: Belong to a type, called with dot notation")
	fmt.Println("• Functions: Standalone, called directly")
	fmt.Println("• Methods: Better encapsulation and OOP-style code")
	fmt.Println("• Functions: More flexible, can work with any parameters")
}

// ========== CHAINING METHODS ==========

// StringBuilder demonstrates method chaining
type StringBuilder struct {
	text string
}

// Append adds text (pointer receiver for chaining)
func (sb *StringBuilder) Append(s string) *StringBuilder {
	sb.text += s
	return sb // Return pointer to allow chaining
}

// AppendLine adds text with newline (pointer receiver for chaining)
func (sb *StringBuilder) AppendLine(s string) *StringBuilder {
	sb.text += s + "\n"
	return sb
}

// Clear clears the content (pointer receiver for chaining)
func (sb *StringBuilder) Clear() *StringBuilder {
	sb.text = ""
	return sb
}

// String returns the final string (value receiver)
func (sb StringBuilder) String() string {
	return sb.text
}

// DemonstrateMethodChaining shows method chaining
func DemonstrateMethodChaining() {
	fmt.Println("\n========== METHOD CHAINING ==========")

	// Method chaining allows multiple method calls in sequence
	sb := &StringBuilder{}
	result := sb.
		AppendLine("Hello, World!").
		AppendLine("This is method chaining.").
		Append("Final line without newline").
		String()

	fmt.Println("Result:")
	fmt.Println(result)

	fmt.Println("\n✓ Method chaining returns the receiver (*StringBuilder)")
	fmt.Println("✓ Allows fluent, readable API design")
}

// RunAllReceiverFunctionDemos runs all receiver function demonstrations
func RunAllReceiverFunctionDemos() {
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║   LEARNING RECEIVER FUNCTIONS         ║")
	fmt.Println("║   (METHODS IN GO)                     ║")
	fmt.Println("╚════════════════════════════════════════╝")

	DemonstrateValueReceivers()
	DemonstratePointerReceivers()
	DemonstrateMethodsOnCustomTypes()
	DemonstrateMethodsWithParameters()
	DemonstrateMethodsVsFunctions()
	DemonstrateMethodChaining()

	fmt.Println("\n" + strings.Repeat("=", 40))
	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("• Value receivers: Work on copies, don't modify original")
	fmt.Println("• Pointer receivers: Work on original, can modify it")
	fmt.Println("• Methods attach behavior to types")
	fmt.Println("• Use pointer receivers for consistency if any method needs to modify")
	fmt.Println(strings.Repeat("=", 40))
}
