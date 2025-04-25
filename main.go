package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

// task with its properties
type Task struct {
	gorm.Model
	ID          string
	Title       string
	Description string
	DueDate     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Status      string
}

func getSingleTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var task Task

	if err := db.First(&task, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task Not Found"})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func getTask(ctx *gin.Context) {
	var tasks []Task
	if err := db.Find(&tasks).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}

func updatedTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var updateData Task

	// Bind only what's present in JSON
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingTask Task
	if err := db.First(&existingTask, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task Not Found"})
		return
	}

	// Only update non-zero/non-empty fields
	if updateData.Title != "" {
		existingTask.Title = updateData.Title
	}
	if updateData.Description != "" {
		existingTask.Description = updateData.Description
	}
	if updateData.Status != "" {
		existingTask.Status = updateData.Status
	}
	if !updateData.DueDate.IsZero() {
		existingTask.DueDate = updateData.DueDate
	}

	if err := db.Save(&existingTask).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task Updated", "task": existingTask})
}

func addTask(ctx *gin.Context) {
	var newTask Task

	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&newTask).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created", "task": newTask})
}

func deleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := db.Delete(&Task{}, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task Not Found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task Removed"})
}

var db *gorm.DB
var err error

func main() {

	db, err = gorm.Open(sqlite.Open("task.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to the Database: ", err)
	}

	db.AutoMigrate(&Task{})

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
