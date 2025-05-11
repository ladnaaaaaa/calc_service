package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ladnaaaaaa/calc_service/internal/database"
	"github.com/ladnaaaaaa/calc_service/internal/orchestrator"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	database.Init()

	server := orchestrator.NewServer()

	if err := server.StartGRPCServer(50051); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := server.Start(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
