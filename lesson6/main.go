package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type Task struct {
	ID          uint   `json:"id,omitempty"`
	Title       string `json:"title,omitempty" binding:"required"`
	Description string `json:"description,omitempty"`
	Status      bool   `json:"status"`
	Priority    uint8  `json:"priority,omitempty"`
}

var tasks = make(map[uint]Task)

//type TaskResponse struct {
// Tasks []Task `json:"tasks"`
// total int    `json:"total"`
//}

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
	err = saveTasksToFile(tasks)
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
	}
}

func getAllTasks(c *gin.Context) {
	statusStr, existsStatus := c.GetQuery("status")
	priorityStr, existsPriority := c.GetQuery("priority")
	if !existsStatus && !existsPriority {
		c.Header("Cache-Control", "public, max-age=3600")
		c.JSON(http.StatusOK, tasks)
		return
	}

	status, err := strconv.ParseBool(statusStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filterTasks := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		if task.Status == status && task.Priority == uint8(priority) {
			filterTasks = append(filterTasks, task)
		}
	}
	c.JSON(http.StatusOK, filterTasks)
}

func updateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
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
	err = saveTasksToFile(tasks)
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
	}
}

func deleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, ok := tasks[uint(id)]
	if ok {
		delete(tasks, uint(id))
		c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
		err = saveTasksToFile(tasks)
		if err != nil {
			c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task not found"})
}

func saveTasksToFile(tasks map[uint]Task) error {
	jsonData, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile("tasks.json", jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func loadTasksFromFile() error {
	data, err := os.ReadFile("tasks.json")
	if os.IsNotExist(err) {
		_, err = os.Create("tasks.json")
		if err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return err
	}
	return nil
}

func listTasks(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	slice := make([]Task, 0, 10)
	for i, task := range tasks {
		if i > uint((page-1)*10) && i <= uint(page*10) {
			slice = append(slice, task)
		}
	}
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].ID < slice[j].ID
	})
	response := make(map[string]interface{})
	response["tasks"] = slice
	response["total"] = len(tasks)
	c.JSON(http.StatusOK, response)
}

func main() {
	err := loadTasksFromFile()
	if err != nil {
		return
	}
	r := gin.Default()

	r.GET("/all", getAllTasks)
	r.POST("/task", createTask)
	r.PUT("/task/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)
	r.GET("/tasks", listTasks)

	r.Run(":8080")
}
