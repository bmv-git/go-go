package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	port := ":8081"
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!")
	})
	err := r.Run(port)
	if err != nil {
		return
	}
}
