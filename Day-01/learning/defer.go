package learning

import "os"

// simple defer - executes after function returns
func simpleDefer() {
	defer cleanup()
	// do work
}

func cleanup() {
	// cleanup code
}

// defer with file operations
func deferFile() {
	f, err := os.Open("file.txt")
	if err != nil {
		return
	}
	defer f.Close() // ensures file is closed when function returns
	
	// read file
}
 
// multiple defers execute in LIFO order (Last In, First Out)
func multipleDeferLIFO() {
	defer first()
	defer second()
	defer third()
	// execution order: third(), second(), first()
}

func first()  {}
func second() {}
func third()  {}

// defer with arguments (evaluated immediately)
func deferWithArgs() {
	x := 10
	defer printValue(x) // x=10 is captured now
	x = 20
	// prints 10, not 20
}

func printValue(v int) {
	_ = v
}

// defer with anonymous function
func deferAnonymous() {
	x := 10
	defer func() {
		_ = x // captures x by reference
	}()
	x = 20
	// anonymous function sees x=20
}

// defer in loop (each iteration adds a defer)
func deferInLoop() {
	for i := 0; i < 3; i++ {
		defer func(n int) {
			_ = n
		}(i)
	}
	// executes: i=2, i=1, i=0
}

// defer for resource cleanup
func deferResourceCleanup() {
	resource := acquireResource()
	defer releaseResource(resource)
	
	// use resource
	// guaranteed cleanup even if panic occurs
}

func acquireResource() int {
	return 1
}

func releaseResource(r int) {
	_ = r
}

// defer modifying named return value
func deferModifyReturn() (result int) {
	defer func() {
		result++ // modifies return value
	}()
	result = 5
	return // returns 6, not 5
}

// defer in error handling
func deferErrorHandling() (err error) {
	f, err := os.Open("file.txt")
	if err != nil {
		return err
	}
	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr // capture close error if no other error
		}
	}()
	
		// process file
		return nil
	}
