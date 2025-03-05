package main

import "github.com/ladnaaaaaa/calc_service/internal/orchestrator"

func main() {
	server := orchestrator.NewServer()
	server.Start(":8080")
}
