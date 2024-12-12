package main

import (
	"aevum-emporium-be/internal/datasource"
	"aevum-emporium-be/internal/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load the PORT from environment variables, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize the database
	client := datasource.ConnectDB()
	if client == nil {
		log.Fatal("Failed to initialize MongoDB connection")
	}

	// Initialize the Gin router
	router := gin.Default() // Initialize once

	// Set up Gin routes (you need to define setupRoutes)
	routes.SetupRoutes(router)

	// Simple ping route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Run the server
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
