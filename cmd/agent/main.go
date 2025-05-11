package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ladnaaaaaa/calc_service/internal/agent"
)

func main() {
	orchestratorAddr := os.Getenv("ORCHESTRATOR_ADDR")
	if orchestratorAddr == "" {
		orchestratorAddr = "localhost:50051"
	}

	client, err := agent.NewCalculatorClient(orchestratorAddr)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := client.ProcessTask(ctx); err != nil {
				log.Printf("failed to process task: %v", err)
				time.Sleep(time.Second)
			}
		}
	}
}
