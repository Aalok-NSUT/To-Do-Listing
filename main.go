package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

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
		saveTasks(tasks)
		fmt.Printf("Successfully added: %s\n", description)

	case "list":
		fmt.Println("---MY TODO LIST---")
		for _, t := range tasks {
			fmt.Printf("%d. [%s] %s\n", t.ID, t.Status, t.Description)
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
		found := false

		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].UpdateStatus(newStatus)
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
			saveTasks(newTasks)
			fmt.Printf("Task %d deleted.\n", id)
		} else {
			fmt.Printf("Task %d not found.\n", id)
		}

	default:
		fmt.Println("Unknown command. Use 'add' or 'list'.")
	}
}
