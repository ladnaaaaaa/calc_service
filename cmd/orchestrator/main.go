package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ladnaaaaaa/calc_service/internal/database"
	"github.com/ladnaaaaaa/calc_service/internal/orchestrator"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Check JWT_SECRET
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	// Initialize database
	database.Init()

	// Create server
	server := orchestrator.NewServer()

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	if err := server.Start(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
