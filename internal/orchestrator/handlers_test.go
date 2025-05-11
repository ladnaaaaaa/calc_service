package orchestrator

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ladnaaaaaa/calc_service/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *Server {
	server := NewServer()
	server.setupRoutes()
	return server
}

func TestCalculateHandler(t *testing.T) {
	server := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/api/v1/calculate",
		strings.NewReader(`{"expression":"2+3*4"}`))
	req.Header.Set("Content-Type", "application/json")

	server.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTaskHandler(t *testing.T) {
	server := setupRouter()

	// Create a test expression
	expr := &models.Expression{
		Expression: "2+3",
		Status:     models.StatusPending,
		UserID:     1,
	}
	err := server.store.AddExpression(expr)
	assert.NoError(t, err)

	// Create test tasks
	task1 := &models.Task{
		ExpressionID: expr.ID,
		Arg1:         2,
		Arg2:         3,
		Operation:    models.OperationAdd,
		Status:       models.StatusPending,
		OrderNum:     0,
	}
	err = server.store.AddTask(task1)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/internal/task", nil)

	server.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
