package taskmanager

import "fmt"

// Custom error types demonstrating error handling

// TaskNotFoundError is returned when a task cannot be found
type TaskNotFoundError struct {
	TaskID int
}

func (e *TaskNotFoundError) Error() string {
	return fmt.Sprintf("task with ID %d not found", e.TaskID)
}

// UserNotFoundError is returned when a user cannot be found
type UserNotFoundError struct {
	UserID int
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user with ID %d not found", e.UserID)
}

// ValidationError is returned when validation fails
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// DuplicateError is returned when trying to create a duplicate entry
type DuplicateError struct {
	EntityType string
	Identifier string
}

func (e *DuplicateError) Error() string {
	return fmt.Sprintf("duplicate %s with identifier '%s' already exists", e.EntityType, e.Identifier)
}

// PermissionError is returned when user lacks permission
type PermissionError struct {
	UserID int
	Action string
}

func (e *PermissionError) Error() string {
	return fmt.Sprintf("user %d does not have permission to perform action: %s", e.UserID, e.Action)
}

// Helper functions to create errors
func NewTaskNotFoundError(taskID int) error {
	return &TaskNotFoundError{TaskID: taskID}
}

func NewUserNotFoundError(userID int) error {
	return &UserNotFoundError{UserID: userID}
}

func NewValidationError(field, message string) error {
	return &ValidationError{Field: field, Message: message}
}

func NewDuplicateError(entityType, identifier string) error {
	return &DuplicateError{EntityType: entityType, Identifier: identifier}
}

func NewPermissionError(userID int, action string) error {
	return &PermissionError{UserID: userID, Action: action}
}
