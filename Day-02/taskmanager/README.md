# Task Manager Project

A comprehensive CLI-based task management system that demonstrates all key Go concepts: **arrays, slices, structs, loops, interfaces, and error handling**.

## Project Structure

```
Day-02/
├── taskmanager/           # Core business logic package
│   ├── models.go         # Data structures (User, Task, Project)
│   ├── errors.go         # Custom error types
│   └── manager.go        # Manager implementation with interfaces
└── cmd/
    └── taskmanager-demo/ # Demo application
        └── main.go       # Main entry point
```

## Features

### 1. **Structs** - Data Models
- `User`: Represents users with ID, name, email, role, and active status
- `Task`: Tasks with title, description, assignee, priority, status, due date, and tags
- `Project`: Collection of tasks with owner and team members

### 2. **Arrays** - Fixed Priority Levels
- `PriorityLevels`: Fixed-size array `[4]string{"Low", "Medium", "High", "Critical"}`
- Tasks reference priority by array index (0-3)

### 3. **Slices** - Dynamic Collections
- Task list: Dynamic slice of tasks
- User list: Dynamic slice of users
- Tags: Each task has a slice of string tags
- Filtering operations return new slices

### 4. **Loops** - Iteration
- `for...range` loops to iterate through tasks and users
- Filtering tasks by status, priority, and tags
- Finding overdue tasks
- Generating statistics

### 5. **Interfaces** - Abstraction
- `TaskRepository`: Interface for task management operations
- `UserRepository`: Interface for user management operations
- `Manager`: Implements both interfaces
- `fmt.Stringer`: Custom String() methods for pretty printing

### 6. **Error Handling** - Custom Errors
- `TaskNotFoundError`: Task lookup failures
- `UserNotFoundError`: User lookup failures
- `ValidationError`: Input validation failures
- `DuplicateError`: Duplicate entry prevention
- `PermissionError`: Authorization failures

## Core Functionality

### Task Operations
- ✅ Create tasks with validation
- ✅ Get task by ID
- ✅ Update task details
- ✅ Delete tasks
- ✅ List all tasks
- ✅ Filter by status (Pending, In Progress, Completed)
- ✅ Filter by priority level
- ✅ Filter by tags
- ✅ Find overdue tasks
- ✅ Assign tasks to users

### User Operations
- ✅ Create users with validation
- ✅ Get user by ID
- ✅ List all users
- ✅ Get active users only
- ✅ Duplicate email prevention

### Statistics
- ✅ Task counts by status
- ✅ Task counts by priority
- ✅ Overdue task tracking
- ✅ Total tasks count

## Running the Demo

```bash
cd Day-02
go run cmd/taskmanager-demo/main.go
```

## Demo Walkthrough

The demo demonstrates:

1. **Creating Users** - Structs with multiple fields
2. **Error Handling** - Duplicate user detection
3. **Arrays** - Fixed priority level array
4. **Creating Tasks** - Complex structs with slices and references
5. **Validation Errors** - Empty field validation
6. **Listing Tasks** - Interface methods and loops
7. **Filtering by Status** - Slices and loops
8. **Filtering by Priority** - Array indexing
9. **Finding Overdue Tasks** - Date comparisons
10. **Filtering by Tags** - Nested slice iteration
11. **Assignment Errors** - Business logic validation
12. **Reassigning Tasks** - Pointer updates
13. **Updating Status** - Struct modification
14. **Active Users** - Filtering with boolean fields
15. **Statistics** - Maps, loops, and arrays combined
16. **Delete Errors** - Not found error handling
17. **Delete Success** - Slice manipulation

## Key Concepts Illustrated

### Arrays vs Slices
```go
// Array (fixed size)
var PriorityLevels = [4]string{"Low", "Medium", "High", "Critical"}

// Slice (dynamic size)
tasks := make([]Task, 0)
tasks = append(tasks, newTask)
```

### Interface Implementation
```go
type TaskRepository interface {
    AddTask(task Task) error
    GetTask(id int) (*Task, error)
    // ... more methods
}

type Manager struct {
    tasks []Task
}

func (m *Manager) AddTask(task Task) error {
    // Implementation
}
```

### Custom Error Types
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}
```

### Loops and Filtering
```go
func (m *Manager) GetTasksByStatus(status TaskStatus) []Task {
    result := make([]Task, 0)
    for _, task := range m.tasks {
        if task.Status == status {
            result = append(result, task)
        }
    }
    return result
}
```

## Real-World Applications

This task manager demonstrates patterns useful for:
- Project management tools
- Bug tracking systems
- Team collaboration platforms
- Workflow automation
- Resource scheduling
- ToDo applications

## Learning Outcomes

After studying this project, you'll understand:
- ✅ How to structure Go packages
- ✅ When to use arrays vs slices
- ✅ How to implement and use interfaces
- ✅ Best practices for error handling
- ✅ Struct composition and pointer usage
- ✅ Loop patterns for filtering and transformation
- ✅ Time handling in Go
- ✅ String formatting and validation
