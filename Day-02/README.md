# Day-02: Advanced Go Concepts & Task Manager Project

Welcome to Day-02 of my Go learning journey! This day focuses on mastering fundamental Go concepts and building a real-world project that ties everything together.

## What I Learned

### Core Concepts Covered

1. **Arrays** 
   - Fixed-size collections
   - Array declaration and initialization
   - Array iteration
   - Default values

2. **Slices** 
   - Dynamic arrays
   - Slice operations (append, slicing)
   - Make and capacity
   - Slice manipulation

3. **Structs** 
   - Custom data types
   - Struct methods
   - Struct composition
   - Pointer receivers

4. **Loops** 
   - For loops
   - While-like loops
   - Range loops
   - Break and continue

5. **Interfaces** 
   - Interface definition
   - Interface implementation
   - Polymorphism
   - Type assertions

6. **Error Handling** 
   - Custom error types
   - Error wrapping
   - Error checking patterns
   - Validation errors

## Project Structure

```
Day-02/
├── README.md                    # This file
├── go.mod                       # Module definition
│
├── learning/                    # Learning examples
│   ├── arrays.go               # Array concepts
│   ├── slices.go               # Slice operations
│   ├── structs.go              # Struct definitions and methods
│   ├── loops.go                # Loop patterns
│   ├── interfaces.go           # Interface examples
│   └── errors.go               # Error handling
│
├── taskmanager/                 # Task Manager Project
│   ├── README.md               # Project documentation
│   ├── models.go               # Data structures (User, Task, Project)
│   ├── errors.go               # Custom error types
│   └── manager.go              # Business logic & interfaces
│
└── cmd/
    ├── learning-demo/
    │   └── main.go             # Demo of learning concepts
    └── taskmanager-demo/
        └── main.go             # Task Manager application demo
```

## Main Project: Task Manager

A comprehensive CLI-based task management system that demonstrates all learned concepts in action.

### Features

**User Management**
- Create users with roles and status
- Active/Inactive user tracking
- Duplicate email prevention
- User listing and filtering

**Task Management**
- Create, read, update, delete tasks
- Task assignment to users
- Priority levels (Low, Medium, High, Critical)
- Status tracking (Pending, In Progress, Completed, Cancelled)
- Due date management
- Tag-based organization

**Advanced Features**
- Overdue task detection
- Filter tasks by status, priority, or tags
- Task statistics dashboard
- Input validation
- Custom error handling

### Technology Highlights

| Concept | Implementation |
|---------|---------------|
| **Arrays** | `PriorityLevels [4]string` - Fixed priority array |
| **Slices** | Dynamic lists for tasks, users, tags |
| **Structs** | `User`, `Task`, `Project` with 10+ fields |
| **Loops** | Filtering, searching, statistics generation |
| **Interfaces** | `TaskRepository`, `UserRepository` - 11 methods |
| **Errors** | 5 custom error types with validation |

### Architecture

```
┌─────────────────────────────────────┐
│     cmd/taskmanager-demo/main.go    │  ← Entry Point
└──────────────┬──────────────────────┘
               │ imports
               ▼
┌─────────────────────────────────────┐
│       taskmanager package           │
├─────────────────────────────────────┤
│  models.go    │ Data Structures     │
│  errors.go    │ Error Types         │
│  manager.go   │ Business Logic      │
└─────────────────────────────────────┘
```

## Running the Code

### Run Learning Examples

```bash
cd Day-02
go run cmd/learning-demo/main.go
```

This demonstrates all the basic concepts with simple examples.

### Run Task Manager Demo

```bash
cd Day-02
go run cmd/taskmanager-demo/main.go
```

This runs a comprehensive demo showcasing:
- Creating users and tasks
- Error handling scenarios
- Filtering and searching
- Task assignments
- Statistics generation
- CRUD operations

## Key Learnings

### 1. Arrays vs Slices
```go
// Array - Fixed size
var priorities [4]string = [4]string{"Low", "Medium", "High", "Critical"}

// Slice - Dynamic size
tasks := make([]Task, 0)
tasks = append(tasks, newTask)
```

### 2. Struct Methods
```go
type Task struct {
    Title    string
    Priority int
}

func (t Task) GetPriorityName() string {
    return PriorityLevels[t.Priority]
}
```

### 3. Interface Implementation
```go
type TaskRepository interface {
    AddTask(task Task) error
    GetTask(id int) (*Task, error)
}

type Manager struct {
    tasks []Task
}

func (m *Manager) AddTask(task Task) error {
    // Implementation
}
```

### 4. Custom Errors
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field '%s': %s", 
        e.Field, e.Message)
}
```

### 5. Loop Patterns
```go
// Filtering with loops
func GetTasksByStatus(status TaskStatus) []Task {
    result := make([]Task, 0)
    for _, task := range tasks {
        if task.Status == status {
            result = append(result, task)
        }
    }
    return result
}
```

## What This Project Teaches

### Beginner Concepts
- Working with arrays and slices
- Creating and using structs
- Loop patterns and iteration
- Basic error handling

### Intermediate Concepts
- Interface design and implementation
- Custom error types
- Pointer usage with structs
- Package organization
- Method receivers

### Advanced Patterns
- Repository pattern (interfaces)
- Error wrapping and handling
- Data validation
- Filtering and transformation
- Statistics aggregation

## Real-World Applications

This task manager demonstrates patterns used in:
- Project management tools (Jira, Trello)
- Bug tracking systems (GitHub Issues)
- Team collaboration platforms
- Workflow automation tools
- Resource scheduling systems
- ToDo applications

## Demo Output Highlights

The task manager demo showcases **17 different scenarios**:

1. Creating users with struct initialization
2. Duplicate detection with custom errors
3. Array-based priority system
4. Complex task creation with slices
5. Input validation errors
6. Listing with interface methods
7. Status-based filtering
8. Priority-based filtering (array indexing)
9. Overdue task detection
10. Tag-based filtering (nested slices)
11. Business rule validation
12. Task reassignment
13. Status updates
14. Active user filtering
15. Comprehensive statistics
16. Error handling for missing entities
17. Slice manipulation for deletion

## Next Steps

After completing Day-02, I have a solid foundation in:
- Go data structures (arrays, slices, structs)
- Control flow (loops, conditionals)
- Interface-based design
- Error handling best practices
- Package organization

**Coming Next**: Day-03 will likely cover:
- Concurrency (goroutines, channels)
- HTTP servers and REST APIs
- Database integration
- Testing and benchmarking

## Notes

- All code follows Go conventions and best practices
- Custom error types provide meaningful error messages
- Interface design allows for easy testing and mocking
- Comprehensive comments explain each concept
- Demo provides hands-on examples of all features

---

**Date**: January 22, 2026  
**Module**: `day02`  
**Go Version**: 1.25.5

Happy Learning! 
