package orchestrator

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) registerWebRoutes() {
	s.engine.Static("/static", "web/static")

	s.engine.GET("/", s.handleHome)
	s.engine.GET("/expressions", s.handleGetExpressionsRequest)
}

func (s *Server) handleHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Distributed Calculator",
	})
}

func (s *Server) handleGetExpressionsRequest(c *gin.Context) {
	expressions := s.store.GetAllExpressions()
	c.JSON(http.StatusOK, gin.H{
		"expressions": expressions,
	})
}

func (s *Server) handleCalculate(c *gin.Context) {
	var req struct {
		Expression string `json:"expression"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid request"})
		return
	}

	tasks, finalTaskID, err := parseExpression(req.Expression)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
	expr := &Expression{
		ID:          id,
		Status:      "processing",
		Tasks:       tasks,
		FinalTaskID: finalTaskID,
		CreatedAt:   time.Now(),
	}

	s.store.AddExpression(expr)
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (s *Server) handleGetExpressions(c *gin.Context) {
	expressions := s.store.GetAllExpressions()
	response := make([]gin.H, 0, len(expressions))

	for _, expr := range expressions {
		response = append(response, gin.H{
			"id":     expr.ID,
			"status": expr.Status,
			"result": expr.Result,
		})
	}

	c.JSON(http.StatusOK, gin.H{"expressions": response})
}

func (s *Server) handleGetExpression(c *gin.Context) {
	id := c.Param("id")
	expr, exists := s.store.GetExpression(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "expression not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"expression": gin.H{
			"id":     expr.ID,
			"status": expr.Status,
			"result": expr.Result,
		},
	})
}

func (s *Server) handleGetTask(c *gin.Context) {
	var availableTask *Task

	for _, task := range s.store.GetPendingTasks() {
		if s.store.IsTaskReady(task) {
			availableTask = task
			break
		}
	}

	if availableTask == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no tasks available"})
		return
	}

	arg1Task, exists1 := s.store.GetTask(availableTask.Arg1ID)
	arg2Task, exists2 := s.store.GetTask(availableTask.Arg2ID)

	if !exists1 || !exists2 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "dependent tasks not found"})
		return
	}

	availableTask.Status = "processing"
	availableTask.StartedAt = time.Now()
	if err := s.store.UpdateTask(availableTask); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task": gin.H{
			"id":             availableTask.ID,
			"arg1_result":    arg1Task.Result,
			"arg2_result":    arg2Task.Result,
			"operation":      availableTask.Operation,
			"operation_time": s.store.GetOperationTime(availableTask.Operation).Milliseconds(),
		},
	})
}

func (s *Server) handleSubmitTask(c *gin.Context) {
	var req struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid request"})
		return
	}

	task, exists := s.store.GetTask(req.ID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	task.Status = "completed"
	task.Result = req.Result
	task.CompletedAt = time.Now()

	if err := s.store.UpdateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	expr, exists := s.store.GetExpression(task.ExpressionID)
	if exists {
		allCompleted := true
		for _, t := range expr.Tasks {
			if t.Status != "completed" {
				allCompleted = false
				break
			}
		}

		if allCompleted {
			expr.Status = "completed"
			expr.Result = task.Result
			expr.UpdatedAt = time.Now()
			s.store.AddExpression(expr)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "result accepted"})

	log.Printf("Updated task: %+v", task)
	log.Printf("Expression status: %+v", expr)
}
