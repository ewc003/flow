package main

import (
	"context"
	"log"
	"os"
	"server/db"
	"server/db/repository"
	"server/handlers"
	"server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file. Using existing environment variables.")
	}

	// Check for required environment variables
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// Initialize MongoDB connection
	client := db.Connect()
	// Ensure the connection is closed when the app exits
	defer client.Disconnect(context.TODO())

	// Initialize repositories
	breakdownRepo := repository.NewBreakdownRepository(client)
	userRepo := repository.NewUserRepository(client)

	// Initialize handlers
	breakdownHandler := handlers.NewBreakdownHandler(breakdownRepo)
	authHandler := handlers.NewAuthHandler(userRepo)

	// Create a Gin router instance
	router := gin.Default()

	// Add a middleware (e.g., logging)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Public routes
	router.GET("/health", handlers.HealthCheckHandler)
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)

	// Create an authenticated group
	authenticated := router.Group("/")
	authenticated.Use(middleware.AuthMiddleware())
	{
		// User routes
		authenticated.GET("/profile", authHandler.GetProfile)

		// Breakdown routes
		authenticated.GET("/breakdowns", breakdownHandler.GetBreakdowns)
		authenticated.GET("/breakdowns/:id", breakdownHandler.GetBreakdownByID)
		authenticated.POST("/breakdowns", breakdownHandler.CreateBreakdown)
		authenticated.PUT("/breakdowns/:id", breakdownHandler.UpdateBreakdown)
		authenticated.DELETE("/breakdowns/:id", breakdownHandler.DeleteBreakdown)
	}

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
