package orchestrator

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

	server.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetTaskHandler(t *testing.T) {
	server := setupRouter()

	numTask1 := &Task{
		ID:     "num_2",
		Result: 2,
		Status: "completed",
	}
	numTask2 := &Task{
		ID:     "num_3",
		Result: 3,
		Status: "completed",
	}
	mainTask := &Task{
		ID:        "task_1",
		Arg1ID:    "num_2",
		Arg2ID:    "num_3",
		Operation: "+",
		Status:    "pending",
		DependsOn: []string{"num_2", "num_3"},
	}

	server.store.AddTask(numTask1)
	server.store.AddTask(numTask2)
	server.store.AddTask(mainTask)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/internal/task", nil)

	server.engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
