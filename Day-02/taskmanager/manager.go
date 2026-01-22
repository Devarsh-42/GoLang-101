package taskmanager

import (
	"strings"
	"time"
)

// TaskRepository interface defines methods for task management
// Demonstrates interface usage
type TaskRepository interface {
	AddTask(task Task) error
	GetTask(id int) (*Task, error)
	UpdateTask(task Task) error
	DeleteTask(id int) error
	ListTasks() []Task
	GetTasksByStatus(status TaskStatus) []Task
	GetTasksByPriority(priority int) []Task
}

// UserRepository interface defines methods for user management
type UserRepository interface {
	AddUser(user User) error
	GetUser(id int) (*User, error)
	ListUsers() []User
	GetActiveUsers() []User
}

// Manager implements both TaskRepository and UserRepository interfaces
type Manager struct {
	tasks      []Task
	users      []User
	nextTaskID int
	nextUserID int
}

// NewManager creates a new Manager instance
func NewManager() *Manager {
	return &Manager{
		tasks:      make([]Task, 0),
		users:      make([]User, 0),
		nextTaskID: 1,
		nextUserID: 1,
	}
}

// AddTask adds a new task (using slices)
func (m *Manager) AddTask(task Task) error {
	// Validation
	if strings.TrimSpace(task.Title) == "" {
		return NewValidationError("Title", "title cannot be empty")
	}

	if task.Priority < 0 || task.Priority >= len(PriorityLevels) {
		return NewValidationError("Priority", "invalid priority level")
	}

	task.ID = m.nextTaskID
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	// Append to slice
	m.tasks = append(m.tasks, task)
	m.nextTaskID++

	return nil
}

// GetTask retrieves a task by ID (using loops)
func (m *Manager) GetTask(id int) (*Task, error) {
	// Loop through tasks slice
	for i := range m.tasks {
		if m.tasks[i].ID == id {
			return &m.tasks[i], nil
		}
	}
	return nil, NewTaskNotFoundError(id)
}

// UpdateTask updates an existing task
func (m *Manager) UpdateTask(task Task) error {
	// Loop to find task
	for i := range m.tasks {
		if m.tasks[i].ID == task.ID {
			task.UpdatedAt = time.Now()
			m.tasks[i] = task
			return nil
		}
	}
	return NewTaskNotFoundError(task.ID)
}

// DeleteTask removes a task (slice manipulation)
func (m *Manager) DeleteTask(id int) error {
	// Loop to find and remove
	for i, task := range m.tasks {
		if task.ID == id {
			// Remove from slice
			m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
			return nil
		}
	}
	return NewTaskNotFoundError(id)
}

// ListTasks returns all tasks
func (m *Manager) ListTasks() []Task {
	return m.tasks
}

// GetTasksByStatus filters tasks by status (loops and slices)
func (m *Manager) GetTasksByStatus(status TaskStatus) []Task {
	result := make([]Task, 0)

	// Loop through tasks and filter
	for _, task := range m.tasks {
		if task.Status == status {
			result = append(result, task)
		}
	}

	return result
}

// GetTasksByPriority filters tasks by priority (using arrays for priority levels)
func (m *Manager) GetTasksByPriority(priority int) []Task {
	if priority < 0 || priority >= len(PriorityLevels) {
		return []Task{}
	}

	result := make([]Task, 0)

	for _, task := range m.tasks {
		if task.Priority == priority {
			result = append(result, task)
		}
	}

	return result
}

// GetOverdueTasks returns all overdue tasks
func (m *Manager) GetOverdueTasks() []Task {
	result := make([]Task, 0)

	for _, task := range m.tasks {
		if task.IsOverdue() {
			result = append(result, task)
		}
	}

	return result
}

// GetTasksByTag filters tasks by tag
func (m *Manager) GetTasksByTag(tag string) []Task {
	result := make([]Task, 0)

	for _, task := range m.tasks {
		// Loop through task tags
		for _, t := range task.Tags {
			if strings.EqualFold(t, tag) {
				result = append(result, task)
				break
			}
		}
	}

	return result
}

// AddUser adds a new user
func (m *Manager) AddUser(user User) error {
	if strings.TrimSpace(user.Name) == "" {
		return NewValidationError("Name", "name cannot be empty")
	}

	if strings.TrimSpace(user.Email) == "" {
		return NewValidationError("Email", "email cannot be empty")
	}

	// Check for duplicate email
	for _, u := range m.users {
		if u.Email == user.Email {
			return NewDuplicateError("user", user.Email)
		}
	}

	user.ID = m.nextUserID
	m.users = append(m.users, user)
	m.nextUserID++

	return nil
}

// GetUser retrieves a user by ID
func (m *Manager) GetUser(id int) (*User, error) {
	for i := range m.users {
		if m.users[i].ID == id {
			return &m.users[i], nil
		}
	}
	return nil, NewUserNotFoundError(id)
}

// ListUsers returns all users
func (m *Manager) ListUsers() []User {
	return m.users
}

// GetActiveUsers returns only active users
func (m *Manager) GetActiveUsers() []User {
	result := make([]User, 0)

	for _, user := range m.users {
		if user.IsActive {
			result = append(result, user)
		}
	}

	return result
}

// AssignTaskToUser assigns a task to a user
func (m *Manager) AssignTaskToUser(taskID, userID int) error {
	task, err := m.GetTask(taskID)
	if err != nil {
		return err
	}

	user, err := m.GetUser(userID)
	if err != nil {
		return err
	}

	if !user.IsActive {
		return NewValidationError("User", "cannot assign task to inactive user")
	}

	task.AssignedTo = user
	task.UpdatedAt = time.Now()

	return m.UpdateTask(*task)
}

// GetTaskStats returns statistics about tasks (demonstrating loops and arrays)
func (m *Manager) GetTaskStats() map[string]int {
	stats := make(map[string]int)

	// Initialize counters for each priority level using array
	for _, priority := range PriorityLevels {
		stats[priority] = 0
	}

	// Count by priority
	for _, task := range m.tasks {
		if task.Priority >= 0 && task.Priority < len(PriorityLevels) {
			priorityName := PriorityLevels[task.Priority]
			stats[priorityName]++
		}
	}

	// Count by status
	stats["Total"] = len(m.tasks)
	stats["Pending"] = 0
	stats["InProgress"] = 0
	stats["Completed"] = 0
	stats["Overdue"] = 0

	for _, task := range m.tasks {
		switch task.Status {
		case StatusPending:
			stats["Pending"]++
		case StatusInProgress:
			stats["InProgress"]++
		case StatusCompleted:
			stats["Completed"]++
		}

		if task.IsOverdue() {
			stats["Overdue"]++
		}
	}

	return stats
}
