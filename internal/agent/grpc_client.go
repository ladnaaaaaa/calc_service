package agent

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ladnaaaaaa/calc_service/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CalculatorClient struct {
	client api.CalculatorClient
	conn   *grpc.ClientConn
}

func NewCalculatorClient(serverAddr string) (*CalculatorClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	client := api.NewCalculatorClient(conn)
	return &CalculatorClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *CalculatorClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *CalculatorClient) GetTask(ctx context.Context) (*api.GetTaskResponse, error) {
	resp, err := c.client.GetTask(ctx, &api.GetTaskRequest{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CalculatorClient) SubmitResult(ctx context.Context, taskID uint64, result float64) error {
	resp, err := c.client.SubmitResult(ctx, &api.SubmitResultRequest{
		TaskId: taskID,
		Result: result,
	})
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("failed to submit result: %s", resp.Error)
	}
	return nil
}

func (c *CalculatorClient) ProcessTask(ctx context.Context) error {
	task, err := c.GetTask(ctx)
	if err != nil {
		return fmt.Errorf("failed to get task: %v", err)
	}

	if task == nil {
		return nil
	}

	var result float64
	switch task.Operation {
	case "+":
		result = task.Arg1 + task.Arg2
	case "-":
		result = task.Arg1 - task.Arg2
	case "*":
		result = task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			return fmt.Errorf("division by zero")
		}
		result = task.Arg1 / task.Arg2
	default:
		return fmt.Errorf("unknown operation: %s", task.Operation)
	}

	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

	if err := c.SubmitResult(ctx, task.TaskId, result); err != nil {
		return fmt.Errorf("failed to submit result: %v", err)
	}

	log.Printf("Processed task %d: %f %s %f = %f", task.TaskId, task.Arg1, task.Operation, task.Arg2, result)
	return nil
}
