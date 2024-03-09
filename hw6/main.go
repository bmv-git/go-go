package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"strconv"
)

type Task struct {
	ID          string `json:"id,omitempty"` // UUID v.4
	Title       string `json:"title,omitempty" binding:"required"`
	Description string `json:"description,omitempty"`
	Status      bool   `json:"status"`
	Priority    uint8  `json:"priority,omitempty"`
}

var tasks = []Task{}     // срез структур Task
var index map[string]int // [ID] = индекс структуры в срезе
var tasksPerPage = 5     // число задач на страницу для пагинации

func createIndex() {
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
	// генерируем строковый ID
	task.ID = uuid.NewString()
	tasks = append(tasks, task)

	// обновляем индекс
	createIndex()

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

	// проверяем, есть ли индекс для данного id
	i, ok := index[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	// если индекс есть, то обновляем i-ю задачу по запросу
	err := c.BindJSON(&tasks[i])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// отправляем клиенту код завершения 200 и обновленную задачу
	c.JSON(http.StatusOK, tasks[i])

	// записываем все задачи в файл
	err = saveTasksToFile(tasks)
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
	}
}

func deleteTask(c *gin.Context) {
	// проверяем параметр DELETE запроса /tasks:id
	id := c.Param("id")

	// проверяем, есть ли индекс для данного id
	i, ok := index[id]
	if ok {
		// сдвигаем срез для удаления i-й задачи
		tasks = append(tasks[:i], tasks[i+1:]...)

		// обновляем индекс
		createIndex()

		// отправляем ответ клиенту и перезаписыаем файл
		c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
		err := saveTasksToFile(tasks)
		if err != nil {
			c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task not found"})
}

func saveTasksToFile(tasks []Task) error {
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
	// проверка деталей GET запроса /tasks?page=
	// (если querry не указан, то номер страницы = 1)
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// выбираем задачи для "нужной" страницы вывода (с номером page из запроса или 1 по умолчанию)...
	// вопрос в другом - по какому критерию выбирать?
	// предполагается, что задачи сортируются по ID.
	// В нашем случае ID имеет строковый тип и не является индексом
	// (имеет строго случайный порядок - по двум причинам:
	// 1) UUID генерирутся случайным образом
	// 2) элементы карты не сортируются в принципе).

	// Таким образом, для пагинации надо создать индекс (вспомогательную карту
	// с UUID в качестве ключа и целочисленным индексом (i) в качестве значения)...
	// При каждом создании или удалении задачи индекс надо будет обновлять.

	iMin := (page - 1) * tasksPerPage
	if iMin >= len(tasks) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "на странице нет задач"})
	}
	iMax := page * tasksPerPage
	if iMax > len(tasks) {
		iMax = len(tasks)
	}
	c.JSON(http.StatusOK, tasks[iMin:iMax])
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
