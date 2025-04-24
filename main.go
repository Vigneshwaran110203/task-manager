package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// task with its properties
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

// mock tasks data
var tasks = []Task{
	{
		ID:          "1",
		Title:       "Task 1",
		Description: "First Task",
		DueDate:     time.Now(),
		Status:      "pending",
	},
	{
		ID:          "2",
		Title:       "Task 2",
		Description: "Second Task",
		DueDate:     time.Now().AddDate(0, 0, 1),
		Status:      "in progress",
	}, {
		ID:          "3",
		Title:       "Task 3",
		Description: "Third Task",
		DueDate:     time.Now().AddDate(0, 0, 2),
		Status:      "finished",
	},
}

func getSingleTask(ctx *gin.Context) {
	id := ctx.Param("id")
	log.Println(id)

	for _, task := range tasks {
		if task.ID == id {
			ctx.JSON(http.StatusOK, task)
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "Task Not Found",
	})
}

func getTask(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}

func updatedTask(ctx *gin.Context) {
	id := ctx.Param("id")

	var updatedTask Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			if updatedTask.Title != "" {
				tasks[i].Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				tasks[i].Description = updatedTask.Description
			}
			if updatedTask.Status != "" {
				tasks[i].Status = updatedTask.Status
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})
			return
		}
	}
}

func addTask(ctx *gin.Context) {
	var newTask Task

	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks = append(tasks, newTask)
	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}

func deleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	for i, val := range tasks {
		if val.ID == id {
			tasks = append(tasks[:1], tasks[i+1:]...)
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Task Removed",
			})
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}

func main() {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "ping",
		})
	})

	router.GET("/tasks", getTask)
	router.GET("/tasks/:id", getSingleTask)
	router.POST("/tasks", addTask)
	router.PUT("/tasks/:id", updatedTask)
	router.DELETE("/tasks/:id", deleteTask)

	router.Run(":8000")
}
