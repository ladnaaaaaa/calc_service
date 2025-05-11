package orchestrator

import (
	"context"
	"fmt"
	"net"

	pb "github.com/ladnaaaaaa/calc_service/api"
	"github.com/ladnaaaaaa/calc_service/internal/models"
	"google.golang.org/grpc"
)

type CalculatorServer struct {
	pb.UnimplementedCalculatorServer
	store *Store
}

func NewCalculatorServer(store *Store) *CalculatorServer {
	return &CalculatorServer{store: store}
}

func (s *CalculatorServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	tasks, err := s.store.GetPendingTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %v", err)
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("no tasks available")
	}

	var readyTask *models.Task
	for i := range tasks {
		if s.store.IsTaskReady(&tasks[i]) {
			readyTask = &tasks[i]
			break
		}
	}

	if readyTask == nil {
		return nil, fmt.Errorf("no tasks available")
	}

	return &pb.GetTaskResponse{
		TaskId:        uint64(readyTask.ID),
		Arg1:          readyTask.Arg1,
		Arg2:          readyTask.Arg2,
		Operation:     string(readyTask.Operation),
		OperationTime: s.store.GetOperationTime(string(readyTask.Operation)).Milliseconds(),
	}, nil
}

func (s *CalculatorServer) SubmitResult(ctx context.Context, req *pb.SubmitResultRequest) (*pb.SubmitResultResponse, error) {
	task, err := s.store.GetTask(uint(req.TaskId))
	if err != nil {
		return &pb.SubmitResultResponse{
			Success: false,
			Error:   "task not found",
		}, nil
	}

	task.Status = models.StatusCompleted
	task.Result = req.Result

	if err := s.store.UpdateTask(task); err != nil {
		return &pb.SubmitResultResponse{
			Success: false,
			Error:   "failed to update task",
		}, nil
	}

	expr, err := s.store.GetExpression(task.ExpressionID)
	if err != nil {
		return &pb.SubmitResultResponse{
			Success: false,
			Error:   "failed to get expression",
		}, nil
	}

	tasks, err := s.store.GetTasksByExpressionID(expr.ID)
	if err != nil {
		return &pb.SubmitResultResponse{
			Success: false,
			Error:   "failed to get tasks",
		}, nil
	}

	allCompleted := true
	for _, t := range tasks {
		if t.Status != models.StatusCompleted {
			allCompleted = false
			break
		}
	}

	if allCompleted {
		expr.Status = models.StatusCompleted
		expr.Result = task.Result
		if err := s.store.UpdateExpression(expr); err != nil {
			return &pb.SubmitResultResponse{
				Success: false,
				Error:   "failed to update expression",
			}, nil
		}
	}

	return &pb.SubmitResultResponse{
		Success: true,
	}, nil
}

func (s *Server) StartGRPCServer(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServer(grpcServer, NewCalculatorServer(s.store))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Printf("failed to serve gRPC: %v\n", err)
		}
	}()

	return nil
}
