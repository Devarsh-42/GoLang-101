package learning

// var keyword with explicit type
var explicitInt int = 42
var explicitString string = "hello"

// var keyword with type inference
var inferredInt = 100
var inferredBool = true

// var keyword without initialization (zero values)
var zeroInt int            // 0
var zeroString string      // ""
var zeroBool bool          // false
var zeroFloat float64      // 0.0
var zeroPointer *int       // nil
var zeroSlice []int        // nil
var zeroMap map[string]int // nil

// multiple variable declarations (same type)
var x, y, z int = 1, 2, 3

// multiple variable declarations (different types)
var (
	name   string  = "Go"
	age    int     = 15
	active bool    = true
	score  float64 = 99.5
)

// short variable declaration (:= only inside functions)
// cannot be used at package level

// constants with explicit type
const maxConnections int = 100
const timeout float64 = 30.5

// constants with type inference
const apiKey = "secret-key"
const retryCount = 3

// multiple constants
const (
	StatusOK       = 200
	StatusNotFound = 404
	StatusError    = 500
)

// iota for enumerated constants
const (
	Sunday    = iota // 0
	Monday           // 1
	Tuesday          // 2
	Wednesday        // 3
	Thursday         // 4
	Friday           // 5
	Saturday         // 6
)

// iota with expressions
const (
	_  = iota             // 0 (ignored)
	KB = 1 << (10 * iota) // 1 << 10 = 1024
	MB                    // 1 << 20 = 1048576
	GB                    // 1 << 30 = 1073741824
)
