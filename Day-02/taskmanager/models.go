package taskmanager

import (
	"fmt"
	"time"
)

// Priority levels using array indexes
var PriorityLevels = [4]string{"Low", "Medium", "High", "Critical"}

// TaskStatus represents the status of a task
type TaskStatus string

const (
	StatusPending    TaskStatus = "Pending"
	StatusInProgress TaskStatus = "In Progress"
	StatusCompleted  TaskStatus = "Completed"
	StatusCancelled  TaskStatus = "Cancelled"
)

// User represents a user in the system
type User struct {
	ID       int
	Name     string
	Email    string
	Role     string
	IsActive bool
}

// Task represents a task in the system
type Task struct {
	ID          int
	Title       string
	Description string
	AssignedTo  *User
	Priority    int // Index to PriorityLevels array (0-3)
	Status      TaskStatus
	DueDate     time.Time
	Tags        []string // Slice of tags
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Project represents a collection of tasks
type Project struct {
	ID          int
	Name        string
	Description string
	Owner       *User
	Tasks       []Task // Slice of tasks
	Members     []User // Slice of team members
	CreatedAt   time.Time
}

// String method for User (implementing fmt.Stringer interface)
func (u User) String() string {
	status := "Active"
	if !u.IsActive {
		status = "Inactive"
	}
	return fmt.Sprintf("%s (%s) - %s [%s]", u.Name, u.Email, u.Role, status)
}

// String method for Task
func (t Task) String() string {
	assignee := "Unassigned"
	if t.AssignedTo != nil {
		assignee = t.AssignedTo.Name
	}
	priority := "Unknown"
	if t.Priority >= 0 && t.Priority < len(PriorityLevels) {
		priority = PriorityLevels[t.Priority]
	}
	return fmt.Sprintf("[%d] %s - %s (Priority: %s, Status: %s, Assigned: %s)",
		t.ID, t.Title, t.Description, priority, t.Status, assignee)
}

// GetPriorityName returns the priority name from array using index
func (t Task) GetPriorityName() string {
	if t.Priority >= 0 && t.Priority < len(PriorityLevels) {
		return PriorityLevels[t.Priority]
	}
	return "Unknown"
}

// IsOverdue checks if task is overdue
func (t Task) IsOverdue() bool {
	return time.Now().After(t.DueDate) && t.Status != StatusCompleted
}

// String method for Project
func (p Project) String() string {
	return fmt.Sprintf("[Project: %s] Owner: %s, Tasks: %d, Members: %d",
		p.Name, p.Owner.Name, len(p.Tasks), len(p.Members))
}
