package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// define a webPort constant
const webPort = ":8082"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load .env file: %v\n", err)
	}

	// make a gin router
	mux := gin.Default()
	// use cors
	mux.Use(cors.Default())
	// setup route
	mux.GET("/api/hello", func(c *gin.Context) {

	})

	// exec router
	mux.Run(webPort)
}
