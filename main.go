package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)

var allowedStatuses = []string{"todo", "in-progress", "done", "blocked"}

type Task struct {
	ID          int      `json:"id"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	Log         []string `json:"log"`
}

func (t *Task) UpdateStatus(newStatus string) {
	t.Status = newStatus
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	entry := "Moved to " + newStatus + " at " + timestamp
	t.Log = append(t.Log, entry)
}

func getStoragePath() string {
	home, _ := os.UserHomeDir()
	return home + "/.gotodo.json"
}

func reindexTasks(tasks []Task) []Task {
	for i := range tasks {
		tasks[i].ID = i + 1
	}
	return tasks
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	path := getStoragePath()
	return os.WriteFile(path, data, 0644)
}

func loadTasks() ([]Task, error) {
	path := getStoragePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []Task{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func isValidStatus(status string) bool {
	for _, s := range allowedStatuses {
		if strings.ToLower(status) == s {
			return true
		}
	}
	return false
}

func printHelp() {
	fmt.Println("gotodo - A simple CLI task manager")
	fmt.Println("\nUsage:")
	fmt.Println("  gotodo [command] [arguments]")
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  add [text]          Add a new task")
	fmt.Println("  list                List all tasks and their logs")
	fmt.Println("  update [id] [stat]  Update task status (e.g., 'Done')")
	fmt.Println("  delete [id]         Remove a task by ID")
	fmt.Println("  search [query]      Find tasks by keyword")
	fmt.Println("  help                Show this menu")
}

func main() {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: gotodo [add|list]")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide a task description")
			return
		}

		description := os.Args[2]

		newTask := Task{
			ID:          len(tasks) + 1,
			Description: description,
			Status:      "Todo",
		}

		newTask.UpdateStatus("Created")

		tasks = append(tasks, newTask)
		tasks = reindexTasks(tasks)
		saveTasks(tasks)
		fmt.Printf("Successfully added: %s , and tasks are re-indexed successfully.\n", description)

	case "list":
		if len(tasks) == 0 {
			fmt.Println("Your list is empty.")
			return
		}

		fmt.Println("ID  | Status      | Description")
		fmt.Println("----|-------------|------------")
		for _, t := range tasks {
			var statusColor string

			// Match the color to the validated status
			switch strings.ToLower(t.Status) {
			case "todo":
				statusColor = ColorYellow
			case "in-progress":
				statusColor = ColorCyan
			case "done":
				statusColor = ColorGreen
			case "blocked":
				statusColor = ColorRed
			default:
				statusColor = ColorReset
			}

			fmt.Printf("%-3d | %s%-11s%s | %s\n",
				t.ID, statusColor, t.Status, ColorReset, t.Description)

			// Let's keep the logs subtle with Cyan
			for _, logEntry := range t.Log {
				fmt.Printf("    %s-> %s%s\n", ColorCyan, logEntry, ColorReset)
			}
			fmt.Println("----|-------------|------------")
		}

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: gotodo update [id] [new_status]")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid ID. Please Provide a number.")
			return
		}

		newStatus := os.Args[3]
		if !isValidStatus(newStatus) {
			fmt.Printf("Error: '%s' is not a valid status.\n", newStatus)
			fmt.Printf("Allowed: %v\n", allowedStatuses)
			return
		}

		found := false

		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].UpdateStatus(strings.ToTitle(newStatus))
				found = true
				break
			}
		}

		if found {
			saveTasks(tasks)
			fmt.Printf("Task %d updated to: %s\n", id, newStatus)
		} else {
			fmt.Printf("Task %d not found.\n", id)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: gotodo delete [id]")
			return
		}

		id, _ := strconv.Atoi(os.Args[2])

		newTasks := []Task{}
		found := false

		for _, t := range tasks {
			if t.ID == id {
				found = true
				continue
			}
			newTasks = append(newTasks, t)
		}

		if found {
			newTasks = reindexTasks(newTasks)
			saveTasks(newTasks)
			fmt.Printf("Task %d deleted, re-indexed all remaining.\n", id)
		} else {
			fmt.Printf("Task %d not found.\n", id)
		}
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Usage: gotodo search [keyword]")
			return
		}

		query := strings.ToLower(os.Args[2])
		fmt.Printf("Searching for: '%s'...\n", query)
		fmt.Println("-------------------------------")

		foundCount := 0

		for _, t := range tasks {
			if strings.Contains(strings.ToLower(t.Description), query) {
				fmt.Printf("%d. [%s] %s\n", t.ID, t.Status, t.Description)
				foundCount++
			}
		}

		if foundCount == 0 {
			fmt.Println("No tasks found matching that keyword")
		} else {
			fmt.Printf("\nFound %d result(s).\n", foundCount)
		}

	case "widget":
		if len(tasks) == 0 {
			fmt.Println("No tasks!")
			return
		}

		found := false

		for _, t := range tasks {
			if strings.ToLower(t.Status) == "done" {
				continue
			}
			fmt.Printf("[%s] %s\n", t.Status, t.Description)
			found = true
		}

		if !found {
			fmt.Println("All caught up! 🎉")
		}

	case "help":
		printHelp()

	default:
		fmt.Printf("Unknown command: '%s'\n\n.", command)
		printHelp()
	}
}
