package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type Task struct {
	ID          string `json:"id,omitempty"` // UUID v.4
	Title       string `json:"title,omitempty" binding:"required"`
	Description string `json:"description,omitempty"`
	Status      bool   `json:"status"`
	Priority    uint8  `json:"priority,omitempty"`
}

var tasks = make([]Task)
var index = make(map[string]int)
var linesPerPage = 10

func createIndex() error {
	index = make(map[string]int)
	for i, task := range tasks {
		key := task.ID
		index[key] = i
	}
}

func createTask(c *gin.Context) {
	var task Task
	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}
	// генерация строкового ID
	task.ID = uuid.NewString()
	tasks = append(tasks, task)
	// обновляем индекс
	_ = createIndex()
	// отправляем сообщение клиенту
	c.JSON(http.StatusOK, gin.H{"message": "задача создана с номером:" + task.ID})
	// записываем файл
	err = saveTasksToFile(tasks)
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
	}
}

func getAllTasks(c *gin.Context) {
	// проверяем детали запроса /all?status=foo&priority=bar
	statusStr, existsStatus := c.GetQuery("status")
	priorityStr, existsPriority := c.GetQuery("priority")
	// если деталей нет, то возвращаем все записи tasks и кэшируем
	if !existsStatus && !existsPriority {
		// кэшируем (где?)
		c.Header("Cache-Control", "public, max-age=3600")
		c.JSON(http.StatusOK, tasks)
		return
	}
	// преобразовываем тип статуса из string в bool
	status, err := strconv.ParseBool(statusStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// преобразовываем тип приоритета из string в int (потом в uint8)
	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// создаем временный срез для возврата отфильтрованных задач (пустой с макс. текущим объемом)
	filterTasks := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		if task.Status == status && task.Priority == uint8(priority) {
			filterTasks = append(filterTasks, task)
		}
	}
	c.JSON(http.StatusOK, filterTasks)
}

func updateTask(c *gin.Context) {
	// проверяем параметр PUT запроса /task:id
	id := c.Param("id")

	i, ok := index[id]
	task, ok := tasks[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// записываем обновленную задачу в карту
	tasks[id] = task
	// отправляем клиенту код завершения 200 и обновленную задачу
	c.JSON(http.StatusOK, task)
	// записываем все задачи в файл
	err = saveTasksToFile(tasks)
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
	}
}

func deleteTask(c *gin.Context) {
	// проверяем параметр DELETE запроса /tasks:id
	id := c.Param("id")

	_, ok := tasks[id]
	if ok {
		delete(tasks, id)
		c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
		err := saveTasksToFile(tasks)
		if err != nil {
			c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task not found"})
}

func saveTasksToFile(tasks map[string]Task) error {
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
	// проверка деталей GET запроса /tasks?page= (если querry не указан, то номер страницы = 1)
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// вспомогательный срез для постраничного вывода
	slice := make([]Task, 0, linesPerPage)
	// выбираем задачи для "нужной" страницы вывода (с номером page из запроса или 1 по умолчанию)...
	// вопрос в другом - по какому критерию выбирать?
	// предполагается, что задачи сортируются по ID.
	// В нашем случае ID имеет строковый тип и не является индексом (имеет строго случайный порядок - по двум причинам:
	// 1) UUID генерирутся случайным образом
	// 2) элементы карты не сортируются в принципе)
	// Таким образом, для пагинации надо создать индекс (вспомогательную карту с UUID в качестве ключа
	// целочисленным индексом (i) в качестве значения)...
	// При каждом создании или удалении задачи индекс надо будет обновлять.
	for i, task := range tasks {
		if i > uint((page-1)*linesPerPage) && i <= uint(page*linesPerPage) {
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
