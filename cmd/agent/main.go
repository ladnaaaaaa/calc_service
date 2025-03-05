package main

import (
	"github.com/ladnaaaaaa/calc_service/internal/agent"
)

func main() {
	agent := agent.NewAgent("http://localhost:8080", 3)
	agent.Start()
}
