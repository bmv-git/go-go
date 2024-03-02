package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.New()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!")
	})

	r.Run(":8081")
}
