package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Task struct {
	ID          uint   `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      bool   `json:"status,omitempty"`
	Priority    uint8  `json:"priority,omitempty"`
}

func createTask(c *gin.Context) {
	var task Task
	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	task.ID = uint(len(tasks) + 1)
	tasks[task.ID] = task
	c.JSON(http.StatusOK, gin.H{"message": "task created"})
}

func getAllTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}

func updateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, ok := tasks[uint(id)]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	err = c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks[uint(id)] = task

	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, ok := tasks[uint(id)]
	if !ok {
		c.JSON(http.StatusOK, gin.H{"message": "task not found"})
		return
	} else {
		delete(tasks, uint(id))
		c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
	}
}

var tasks = make(map[uint]Task)

func main() {
	r := gin.Default()

	r.GET("/all", getAllTasks)
	r.POST("/task", createTask)
	r.PUT("/task/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)

	r.Run(":8080")
}
