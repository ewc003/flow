package main

import (
	"context"
	"server/db"
	"server/db/repository"
	"server/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize MongoDB connection
	client := db.Connect()
	// Ensure the connection is closed when the app exits
	defer client.Disconnect(context.TODO())

	// Initialize repositories
	breakdownRepo := repository.NewBreakdownRepository(client)

	// Initialize handlers
	breakdownHandler := handlers.NewBreakdownHandler(breakdownRepo)

	// Create a Gin router instance
	router := gin.Default()

	// Add a middleware (e.g., logging)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Define routes
	router.GET("/health", handlers.HealthCheckHandler)
	router.POST("/breakdowns", breakdownHandler.CreateBreakdowns)
	router.GET("/breakdowns", breakdownHandler.GetBreakdowns)

	// Start the server
	port := ":8080"
	router.Run(port)
}
