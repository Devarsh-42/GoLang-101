package main

import (
	"fmt"
	"time"

	"day02/taskmanager"
)

func main() {
	fmt.Println("=== Task Manager Demo ===\n")

	// Create a new manager
	manager := taskmanager.NewManager()

	// 1. STRUCTS: Create users
	fmt.Println("1. Creating Users (Structs):")
	users := []taskmanager.User{
		{
			Name:     "Alice Johnson",
			Email:    "alice@example.com",
			Role:     "Developer",
			IsActive: true,
		},
		{
			Name:     "Bob Smith",
			Email:    "bob@example.com",
			Role:     "Manager",
			IsActive: true,
		},
		{
			Name:     "Charlie Brown",
			Email:    "charlie@example.com",
			Role:     "Designer",
			IsActive: false, // Inactive user
		},
	}

	// SLICES & LOOPS: Add users
	for _, user := range users {
		if err := manager.AddUser(user); err != nil {
			fmt.Printf("   Error adding user: %v\n", err)
		} else {
			fmt.Printf("   ✓ Added user: %s\n", user.Name)
		}
	}
	fmt.Println()

	// 2. ERRORS: Try to add duplicate user
	fmt.Println("2. Error Handling - Adding Duplicate User:")
	duplicateUser := taskmanager.User{
		Name:     "Alice Duplicate",
		Email:    "alice@example.com", // Duplicate email
		Role:     "Tester",
		IsActive: true,
	}
	if err := manager.AddUser(duplicateUser); err != nil {
		fmt.Printf("   ✗ Error (expected): %v\n\n", err)
	}

	// 3. ARRAYS: Display priority levels
	fmt.Println("3. Arrays - Priority Levels:")
	fmt.Printf("   Available priority levels: ")
	for i, priority := range taskmanager.PriorityLevels {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%s (%d)", priority, i)
	}
	fmt.Println("\n")

	// 4. Create tasks with different priorities
	fmt.Println("4. Creating Tasks (Structs, Slices, Arrays):")

	// Get users for assignment
	alice, _ := manager.GetUser(1)
	bob, _ := manager.GetUser(2)

	tasks := []taskmanager.Task{
		{
			Title:       "Implement User Authentication",
			Description: "Add JWT-based authentication",
			AssignedTo:  alice,
			Priority:    3, // Critical
			Status:      taskmanager.StatusInProgress,
			DueDate:     time.Now().Add(48 * time.Hour),
			Tags:        []string{"security", "backend", "urgent"},
		},
		{
			Title:       "Design Landing Page",
			Description: "Create mockups for new landing page",
			AssignedTo:  nil, // Unassigned
			Priority:    1,   // Medium
			Status:      taskmanager.StatusPending,
			DueDate:     time.Now().Add(5 * 24 * time.Hour),
			Tags:        []string{"design", "frontend"},
		},
		{
			Title:       "Fix Login Bug",
			Description: "Users can't login with special characters",
			AssignedTo:  alice,
			Priority:    2, // High
			Status:      taskmanager.StatusPending,
			DueDate:     time.Now().Add(-24 * time.Hour), // Overdue!
			Tags:        []string{"bug", "urgent", "backend"},
		},
		{
			Title:       "Update Documentation",
			Description: "Add API documentation for new endpoints",
			AssignedTo:  bob,
			Priority:    0, // Low
			Status:      taskmanager.StatusCompleted,
			DueDate:     time.Now().Add(7 * 24 * time.Hour),
			Tags:        []string{"documentation"},
		},
		{
			Title:       "Code Review",
			Description: "Review pull requests from team",
			AssignedTo:  bob,
			Priority:    1, // Medium
			Status:      taskmanager.StatusInProgress,
			DueDate:     time.Now().Add(24 * time.Hour),
			Tags:        []string{"review", "management"},
		},
	}

	// LOOPS & ERROR HANDLING: Add tasks
	for _, task := range tasks {
		if err := manager.AddTask(task); err != nil {
			fmt.Printf("   Error adding task: %v\n", err)
		} else {
			fmt.Printf("   ✓ Added: %s [%s]\n", task.Title, task.GetPriorityName())
		}
	}
	fmt.Println()

	// 5. VALIDATION ERROR: Try to add invalid task
	fmt.Println("5. Error Handling - Validation:")
	invalidTask := taskmanager.Task{
		Title:    "", // Empty title
		Priority: 1,
		Status:   taskmanager.StatusPending,
		DueDate:  time.Now(),
	}
	if err := manager.AddTask(invalidTask); err != nil {
		fmt.Printf("   ✗ Validation error (expected): %v\n\n", err)
	}

	// 6. INTERFACES & LOOPS: List all tasks
	fmt.Println("6. Listing All Tasks (Interface Methods, Loops):")
	allTasks := manager.ListTasks()
	for _, task := range allTasks {
		fmt.Printf("   %s\n", task.String())
	}
	fmt.Println()

	// 7. SLICES & LOOPS: Filter by status
	fmt.Println("7. Filtering Tasks by Status (Slices, Loops):")
	pendingTasks := manager.GetTasksByStatus(taskmanager.StatusPending)
	fmt.Printf("   Pending tasks (%d):\n", len(pendingTasks))
	for _, task := range pendingTasks {
		fmt.Printf("   - %s (Priority: %s)\n", task.Title, task.GetPriorityName())
	}
	fmt.Println()

	// 8. ARRAYS & LOOPS: Filter by priority
	fmt.Println("8. Filtering Tasks by Priority (Arrays, Loops):")
	criticalPriority := 3 // Index to PriorityLevels array
	criticalTasks := manager.GetTasksByPriority(criticalPriority)
	fmt.Printf("   %s priority tasks (%d):\n", taskmanager.PriorityLevels[criticalPriority], len(criticalTasks))
	for _, task := range criticalTasks {
		fmt.Printf("   - %s\n", task.Title)
	}
	fmt.Println()

	// 9. Check overdue tasks
	fmt.Println("9. Finding Overdue Tasks (Loops, Slices):")
	overdueTasks := manager.GetOverdueTasks()
	fmt.Printf("   Overdue tasks (%d):\n", len(overdueTasks))
	for _, task := range overdueTasks {
		fmt.Printf("   ⚠ %s (Due: %s)\n", task.Title, task.DueDate.Format("2006-01-02"))
	}
	fmt.Println()

	// 10. SLICES: Filter by tag
	fmt.Println("10. Filtering Tasks by Tag (Slices, Loops):")
	urgentTasks := manager.GetTasksByTag("urgent")
	fmt.Printf("   Tasks tagged 'urgent' (%d):\n", len(urgentTasks))
	for _, task := range urgentTasks {
		tags := ""
		for i, tag := range task.Tags {
			if i > 0 {
				tags += ", "
			}
			tags += tag
		}
		fmt.Printf("   - %s [Tags: %s]\n", task.Title, tags)
	}
	fmt.Println()

	// 11. INTERFACES: Try to assign task to inactive user (error handling)
	fmt.Println("11. Error Handling - Assigning to Inactive User:")
	taskToAssign, _ := manager.GetTask(2)
	inactiveUser, _ := manager.GetUser(3)
	fmt.Printf("   Attempting to assign '%s' to %s (inactive)...\n", taskToAssign.Title, inactiveUser.Name)
	if err := manager.AssignTaskToUser(2, 3); err != nil {
		fmt.Printf("   ✗ Error (expected): %v\n\n", err)
	}

	// 12. Successfully reassign task
	fmt.Println("12. Reassigning Task (Interfaces):")
	fmt.Printf("   Assigning '%s' to %s...\n", taskToAssign.Title, bob.Name)
	if err := manager.AssignTaskToUser(2, 2); err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		updatedTask, _ := manager.GetTask(2)
		fmt.Printf("   ✓ Task assigned to: %s\n\n", updatedTask.AssignedTo.Name)
	}

	// 13. Update task status
	fmt.Println("13. Updating Task Status (Structs):")
	taskToUpdate, _ := manager.GetTask(3)
	fmt.Printf("   Updating task '%s' from %s to %s...\n",
		taskToUpdate.Title, taskToUpdate.Status, taskmanager.StatusCompleted)
	taskToUpdate.Status = taskmanager.StatusCompleted
	if err := manager.UpdateTask(*taskToUpdate); err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Task updated successfully\n\n")
	}

	// 14. Get active users
	fmt.Println("14. Listing Active Users (Slices, Loops):")
	activeUsers := manager.GetActiveUsers()
	fmt.Printf("   Active users (%d):\n", len(activeUsers))
	for _, user := range activeUsers {
		fmt.Printf("   - %s\n", user.String())
	}
	fmt.Println()

	// 15. ARRAYS & LOOPS: Display statistics
	fmt.Println("15. Task Statistics (Arrays, Loops, Maps):")
	stats := manager.GetTaskStats()
	fmt.Printf("   Total Tasks: %d\n", stats["Total"])
	fmt.Printf("   Status Breakdown:\n")
	fmt.Printf("     - Pending: %d\n", stats["Pending"])
	fmt.Printf("     - In Progress: %d\n", stats["InProgress"])
	fmt.Printf("     - Completed: %d\n", stats["Completed"])
	fmt.Printf("     - Overdue: %d\n", stats["Overdue"])
	fmt.Printf("   Priority Breakdown:\n")
	for i, priority := range taskmanager.PriorityLevels {
		fmt.Printf("     - %s: %d\n", priority, stats[taskmanager.PriorityLevels[i]])
	}
	fmt.Println()

	// 16. ERROR HANDLING: Try to delete non-existent task
	fmt.Println("16. Error Handling - Deleting Non-existent Task:")
	if err := manager.DeleteTask(999); err != nil {
		fmt.Printf("   ✗ Error (expected): %v\n\n", err)
	}

	// 17. Successfully delete a task
	fmt.Println("17. Deleting a Task (Slices manipulation):")
	taskToDelete, _ := manager.GetTask(4)
	fmt.Printf("   Deleting task: '%s'...\n", taskToDelete.Title)
	if err := manager.DeleteTask(4); err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Task deleted successfully\n")
		fmt.Printf("   Remaining tasks: %d\n\n", len(manager.ListTasks()))
	}

	fmt.Println("=== Demo Complete ===")
	fmt.Println("\n✓ Demonstrated Concepts:")
	fmt.Println("  • Structs: User, Task with multiple fields")
	fmt.Println("  • Arrays: PriorityLevels fixed-size array")
	fmt.Println("  • Slices: Dynamic lists of tasks, users, tags")
	fmt.Println("  • Loops: Iterating over slices with for...range")
	fmt.Println("  • Interfaces: TaskRepository, UserRepository implementations")
	fmt.Println("  • Errors: Custom error types and validation")
}
