package main

import (
	"fmt"
	"log"

	"go-chat-app/config"
	"go-chat-app/handlers"
	"go-chat-app/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	// Connect to database
	err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Create a new router
	router := gin.Default()

	// Auth Routes
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	// Protected WebSocket Route
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// WebSocket route (real-time messaging)
	protected.GET("/ws", handlers.HandleWebSocket)

	// REST API route (chat history with pagination)
	protected.GET("/chat/history", handlers.FetchChatHistory)

	// Start server
	port := "8080"
	fmt.Println("Server running on port " + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
