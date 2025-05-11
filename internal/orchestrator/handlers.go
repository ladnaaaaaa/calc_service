package orchestrator

import (
	"net/http"
	"strconv"

	"github.com/ladnaaaaaa/calc_service/internal/models"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerWebRoutes() {
	s.Engine.Static("/static", "web/static")

	s.Engine.GET("/", s.handleHome)
	s.Engine.GET("/expressions", s.handleGetExpressionsRequest)
}

func (s *Server) handleHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Distributed Calculator",
	})
}

func (s *Server) handleGetExpressionsRequest(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	expressions, err := s.store.GetAllExpressions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expressions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"expressions": expressions,
	})
}

func (s *Server) handleCalculate(c *gin.Context) {
	var req struct {
		Expression string `json:"expression" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tasks, err := ParseExpression(req.Expression)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expr := &models.Expression{
		Expression: req.Expression,
		Status:     models.StatusPending,
		UserID:     userID,
	}

	if err := s.store.AddExpression(expr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save expression"})
		return
	}

	for _, task := range tasks {
		task.ExpressionID = expr.ID
		if err := s.store.AddTask(task); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save task"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     expr.ID,
		"status": expr.Status,
	})
}

func (s *Server) handleGetExpressions(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	expressions, err := s.store.GetAllExpressions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expressions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expressions": expressions})
}

func (s *Server) handleGetExpression(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	expr, err := s.store.GetExpression(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "expression not found"})
		return
	}

	if expr.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, expr)
}

func (s *Server) handleGetTask(c *gin.Context) {
	tasks, err := s.store.GetPendingTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tasks"})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no tasks available"})
		return
	}

	// Find the first task that is ready to be executed
	var readyTask *models.Task
	for i := range tasks {
		if s.store.IsTaskReady(&tasks[i]) {
			readyTask = &tasks[i]
			break
		}
	}

	if readyTask == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no tasks available"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task": gin.H{
			"id":             readyTask.ID,
			"arg1_result":    readyTask.Arg1,
			"arg2_result":    readyTask.Arg2,
			"operation":      readyTask.Operation,
			"operation_time": s.store.GetOperationTime(string(readyTask.Operation)).Milliseconds(),
		},
	})
}

func (s *Server) handleSubmitTask(c *gin.Context) {
	var req struct {
		ID     uint    `json:"id" binding:"required"`
		Result float64 `json:"result" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	task, err := s.store.GetTask(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	task.Status = models.StatusCompleted
	task.Result = req.Result

	if err := s.store.UpdateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	expr, err := s.store.GetExpression(task.ExpressionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expression"})
		return
	}

	tasks, err := s.store.GetTasksByExpressionID(expr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tasks"})
		return
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update expression"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "result accepted"})
}
