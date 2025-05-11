package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ladnaaaaaa/calc_service/internal/database"
	"github.com/ladnaaaaaa/calc_service/internal/orchestrator"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *orchestrator.Server {
	// Set test environment variables
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("TIME_ADDITION_MS", "100")
	os.Setenv("TIME_SUBTRACTION_MS", "100")
	os.Setenv("TIME_MULTIPLICATIONS_MS", "100")
	os.Setenv("TIME_DIVISIONS_MS", "100")
	os.Setenv("DB_PATH", "test.db")

	// Remove test DB before each run
	_ = os.Remove("test.db")

	// Initialize database
	database.Init()

	// Create server
	return orchestrator.NewServer()
}

func TestIntegration(t *testing.T) {
	server := setupTestServer()

	// Test registration
	t.Run("Registration", func(t *testing.T) {
		reqBody := map[string]string{
			"login":    "testuser",
			"password": "testpass",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		server.Engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test login
	var token string
	t.Run("Login", func(t *testing.T) {
		reqBody := map[string]string{
			"login":    "testuser",
			"password": "testpass",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		server.Engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		token = response["token"]
		assert.NotEmpty(t, token)
	})

	// Test expression calculation
	var expressionID uint
	t.Run("Calculate Expression", func(t *testing.T) {
		reqBody := map[string]string{
			"expression": "2+3*4",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		server.Engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		expressionID = uint(response["id"].(float64))
		assert.NotZero(t, expressionID)
	})

	// Test get expressions list
	t.Run("Get Expressions", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/expressions", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		server.Engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		expressions := response["expressions"].([]interface{})
		assert.NotEmpty(t, expressions)
	})

	// Test get specific expression
	t.Run("Get Expression", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/expressions/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		server.Engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		expr := response["expression"].(map[string]interface{})
		assert.Equal(t, "2+3*4", expr["Expression"])
	})

	// Test task processing
	t.Run("Task Processing", func(t *testing.T) {
		// Get task
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/internal/task", nil)
		server.Engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var taskResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &taskResponse)
		task := taskResponse["task"].(map[string]interface{})
		taskID := task["id"].(float64)

		// Submit task result
		reqBody := map[string]interface{}{
			"id":     taskID,
			"result": 14.0,
		}
		jsonBody, _ := json.Marshal(reqBody)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/internal/task", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		server.Engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Wait for expression to complete
		time.Sleep(200 * time.Millisecond)

		// Check expression status
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/expressions/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		server.Engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var exprResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &exprResponse)
		expr := exprResponse["expression"].(map[string]interface{})
		assert.Equal(t, "completed", expr["Status"])
		assert.Equal(t, 14.0, expr["Result"])
	})
}

// Добавляю edge-case тесты
func TestIntegration_EdgeCases(t *testing.T) {
	server := setupTestServer()

	t.Run("Invalid Token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/expressions", nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")
		server.Engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Repeat Registration", func(t *testing.T) {
		body := map[string]string{"login": "user1", "password": "pass1"}
		jsonBody, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		server.Engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonBody))
		req2.Header.Set("Content-Type", "application/json")
		server.Engine.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusBadRequest, w2.Code)
	})

	t.Run("Nonexistent Expression", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/expressions/9999", nil)
		server.Engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code) // Без токена
	})

	t.Run("Nonexistent Task", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := map[string]interface{}{"id": 9999, "result": 42.0}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/internal/task", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		server.Engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
