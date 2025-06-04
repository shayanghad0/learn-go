package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "strconv"
)

const dataFile = "tasks.json"

type Task struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Done   bool   `json:"done"`
}

var tasks []Task

func loadTasks() {
    data, err := ioutil.ReadFile(dataFile)
    if err == nil {
        json.Unmarshal(data, &tasks)
    }
}

func saveTasks() {
    data, _ := json.MarshalIndent(tasks, "", "  ")
    ioutil.WriteFile(dataFile, data, 0644)
}

func addTask(title string) {
    id := 1
    if len(tasks) > 0 {
        id = tasks[len(tasks)-1].ID + 1
    }
    task := Task{ID: id, Title: title, Done: false}
    tasks = append(tasks, task)
    saveTasks()
    fmt.Println("Added:", title)
}

func listTasks() {
    for _, task := range tasks {
        status := "❌"
        if task.Done {
            status = "✅"
        }
        fmt.Printf("[%d] %s %s\n", task.ID, status, task.Title)
    }
}

func markDone(id int) {
    for i, task := range tasks {
        if task.ID == id {
            tasks[i].Done = true
            saveTasks()
            fmt.Println("Marked done:", task.Title)
            return
        }
    }
    fmt.Println("Task not found.")
}

func deleteTask(id int) {
    for i, task := range tasks {
        if task.ID == id {
            tasks = append(tasks[:i], tasks[i+1:]...)
            saveTasks()
            fmt.Println("Deleted:", task.Title)
            return
        }
    }
    fmt.Println("Task not found.")
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: todo add|list|done|delete <arg>")
        return
    }

    loadTasks()

    switch os.Args[1] {
    case "add":
        if len(os.Args) < 3 {
            fmt.Println("Title is required.")
            return
        }
        addTask(os.Args[2])
    case "list":
        listTasks()
    case "done":
        if len(os.Args) < 3 {
            fmt.Println("ID is required.")
            return
        }
        id, _ := strconv.Atoi(os.Args[2])
        markDone(id)
    case "delete":
        if len(os.Args) < 3 {
            fmt.Println("ID is required.")
            return
        }
        id, _ := strconv.Atoi(os.Args[2])
        deleteTask(id)
    default:
        fmt.Println("Unknown command.")
    }
}
