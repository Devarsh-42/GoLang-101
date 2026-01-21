# Go Learning - Day 01

Learning Go fundamentals with idiomatic practices and clean code principles.

## Structure

```
.
├── cmd/                    # CLI entry points
│   └── calculator/        # Calculator CLI application
├── internal/              # Private application logic
│   └── calculator/        # Calculator business logic and tests
├── learning/              # Go syntax reference examples
│   ├── variables.go       # Variable and constant declarations
│   ├── control_flow.go    # if, for, switch, type switch
│   ├── functions.go       # Function syntax and patterns
│   └── defer.go          # Defer statement examples
├── go.mod                 # Go module definition
└── README.md             # This file
```

## Learning Files

The `learning/` directory contains reference examples:
- **variables.go** - var vs :=, zero values, constants, iota
- **control_flow.go** - if, for loops, switch, type switch
- **functions.go** - parameters, returns, named returns, methods
- **defer.go** - defer execution order, resource cleanup patterns

These files compile but are for reference only.

## Calculator Application

A simple calculator demonstrating Go project structure.

### Run Tests

```bash
go test ./internal/calculator/...
```

### Run Application

```bash
go run cmd/calculator/main.go
```

## Principles Followed

- No global variables
- Explicit error handling
- Thin main() function
- Testable business logic
- Standard library only
- Idiomatic Go code
