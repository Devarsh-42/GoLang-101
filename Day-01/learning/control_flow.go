package learning

import "fmt"

// if statement
func exampleIf() {
	age := 18
	if age >= 18 {
		// can vote
	}
}

// if-else
func ExampleIfElse() {
	score := 75
	if score >= 60 {
		fmt.Println("THE STUDENT IS PASSEDDD!!!!!")
	} else {
		fmt.Println("THE STUDENT HAS FAILED!!!!!")
	}
}

// if-else-if
func exampleIfElseIf() {
	grade := 85
	if grade >= 90 {
		// A
	} else if grade >= 80 {
		// B
	} else if grade >= 70 {
		// C
	} else {
		// F
	}
}

// if with short statement
func exampleIfShortStatement() {
	// variable scope limited to if block
	if num := compute(); num > 10 {
		// num is available here
	} else {
		// and here
	}
	// but not here
}

func compute() int {
	return 42
}

// for loop - classic style
func exampleForClassic() {
	for i := 0; i < 10; i++ {
		// loop body
	}
}

// for loop - while style (only condition)
func exampleForWhile() {
	count := 0
	for count < 5 {
		count++
	}
}

// for loop - infinite loop
func exampleForInfinite() {
	for {
		// infinite loop
		// break to exit
		break
	}
}

// for loop - range over slice
func exampleForRangeSlice() {
	numbers := []int{1, 2, 3, 4, 5}

	// with index and value
	for i, v := range numbers {
		_, _ = i, v
	}

	// only index
	for i := range numbers {
		_ = i
	}

	// only value (use _ to ignore index)
	for _, v := range numbers {
		_ = v
	}
}

// for loop - range over map
func exampleForRangeMap() {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	for key, value := range m {
		_, _ = key, value
	}
}

// for loop - range over string (runes)
func exampleForRangeString() {
	s := "hello"

	for i, r := range s {
		_, _ = i, r // i is byte index, r is rune
	}
}

// switch - expression
func exampleSwitch() {
	day := "Monday"

	switch day {
	case "Monday":
		// do Monday things
	case "Tuesday", "Wednesday":
		// multiple cases
	case "Friday":
		// do Friday things
	default:
		// default case
	}
}

// switch - no expression (like if-else-if)
func exampleSwitchNoExpression() {
	age := 25

	switch {
	case age < 13:
		// child
	case age < 20:
		// teenager
	case age < 60:
		// adult
	default:
		// senior
	}
}

// switch - with short statement
func exampleSwitchShortStatement() {
	switch num := compute(); {
	case num < 0:
		// negative
	case num == 0:
		// zero
	default:
		// positive
	}
}

// type switch
func exampleTypeSwitch() {
	var i interface{} = "hello"

	switch v := i.(type) {
	case int:
		_ = v // v is int
	case string:
		_ = v // v is string
	case bool:
		_ = v // v is bool
	default:
		_ = v // v is interface{}
	}
}

// continue and break
func exampleContinueBreak() {
	for i := 0; i < 10; i++ {
		if i == 5 {
			continue // skip this iteration
		}
		if i == 8 {
			break // exit loop
		}
	}
}
