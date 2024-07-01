package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const webPort = ":8082"

func main() {
	mux := gin.Default()

	mux.Use(cors.Default())

	mux.GET("/api/hello", func(c *gin.Context) {

	})

	mux.Run(webPort)
}
