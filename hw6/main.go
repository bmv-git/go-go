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

var tasks []Task         // срез структур Task
var index map[string]int // [ID] = индекс структуры в срезе
var tasksPerPage = 5     // число задач на страницу для пагинации

func createIndex() {
	index = make(map[string]int)
	for i, task := range tasks {
		index[task.ID] = i
	}
}

// обработчик запроса POST /task
func createTask(c *gin.Context) {
	// создаем новую задачу
	var task Task

	//привязываем данные запроса (JSON) к задаче task (структуры Тask)
	// *) BindJSON - это реализация интерфейса Binding для входных данных
	// в формате JSON.

	// type Binding interface {
	//	Name() string
	//	Bind(*http.Request, any) error
	// }
	// Binding describes the interface which needs to be implemented
	// for binding the data present in the request such as JSON request body,
	// query parameters or the form POST.
	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// генерируем строковый ID (UUID v.4)
	task.ID = uuid.NewString()

	// записываем задачу в срез задач
	tasks = append(tasks, task)

	// обновляем индекс
	createIndex()

	// отправляем сообщение клиенту
	// func (c *Context) JSON(code int, obj any)
	// JSON serializes the given struct as JSON into the response body.
	// It also sets the Content-Type as "application/json".

	// type H map[string]any
	// H is a shortcut for map[string]any
	c.JSON(http.StatusOK, gin.H{"message": "задача создана с номером: " + task.ID})

	// записываем файл
	err = saveTasksToFile(tasks)
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
	}
}

// обработчик запроса GET /all?status=  &priority=
func getAllTasks(c *gin.Context) {
	// проверяем детали запроса
	statusStr, existsStatus := c.GetQuery("status")
	priorityStr, existsPriority := c.GetQuery("priority")

	// если деталей нет, то возвращаем все записи tasks и кэшируем
	if !existsStatus && !existsPriority {

		// кэшируем (где?)
		c.Header("Cache-Control", "public, max-age=3600")
		// возвращаем все записи tasks
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

// обработчик запроса PUT /task/:id
func updateTask(c *gin.Context) {
	// проверяем параметр
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
		return
	}
}

// обработчик запроса DELETE /tasks/:id
func deleteTask(c *gin.Context) {
	// проверяем параметр
	id := c.Param("id")

	// проверяем, есть ли индекс для данного id
	i, ok := index[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	// сдвигаем срез для удаления i-й задачи
	tasks = append(tasks[:i], tasks[i+1:]...)

	// обновляем индекс
	createIndex()

	// отправляем ответ клиенту и перезаписыаем файл
	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})

	// записываем срез задач в файл
	err := saveTasksToFile(tasks)
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
		return
	}
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

// обработчик запроса GET /tasks?page=
func listTasks(c *gin.Context) {
	// (если ?query не указан, то номер страницы = 1)
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
	// 1) UUID генерируются случайным образом
	// 2) элементы карты не сортируются в принципе).

	// Таким образом, для пагинации надо создать индекс (вспомогательную карту
	// с UUID в качестве ключа и целочисленным индексом (i) в качестве значения)...
	// При каждом создании или удалении задачи индекс надо будет обновлять.

	// страница формируется сразу из среза задач по нижнему и верхнему индексу
	iL := (page - 1) * tasksPerPage
	if iL >= len(tasks) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "на этой странице нет задач"})
		return
	}
	iH := page * tasksPerPage
	if iH > len(tasks) {
		iH = len(tasks)
	}
	c.JSON(http.StatusOK, tasks[iL:iH])
}
func homePage(c *gin.Context) {
	c.String(http.StatusOK, "СПИСОК ЗАДАЧ\n")
}

func main() {
	err := loadTasksFromFile()
	if err != nil {
		return
	}

	// обновляем индекс
	createIndex()

	r := gin.Default()

	r.GET("/", homePage)
	r.POST("/task", createTask)
	r.GET("/all", getAllTasks)
	r.GET("/tasks", listTasks)
	r.PUT("/task/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)

	err = r.Run(":8080")
	if err != nil {
		return
	}
}
