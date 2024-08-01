package main

import (
	"fmt"

	"github.com/anakilang-ai/backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		routes.URL(c.Writer, c.Request)
	})
	port := ":8080"
	fmt.Println("Server started at: http://localhost" + port)
	r.Run(port)
}