package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ladnaaaaaa/calc_service/internal/agent"
	"github.com/stretchr/testify/assert"
)

func setupTestAgent() *agent.Agent {
	// Set test environment variables
	os.Setenv("ORCHESTRATOR_URL", "http://localhost:8080")
	os.Setenv("COMPUTING_POWER", "2")
	os.Setenv("TIME_ADDITION_MS", "100")
	os.Setenv("TIME_SUBTRACTION_MS", "100")
	os.Setenv("TIME_MULTIPLICATIONS_MS", "100")
	os.Setenv("TIME_DIVISIONS_MS", "100")

	// Create agent
	return agent.NewAgent("http://localhost:8080", 2)
}

func TestAgent(t *testing.T) {
	agent := setupTestAgent()

	// Start agent
	go agent.Start()

	// Wait for agent to start
	time.Sleep(100 * time.Millisecond)

	// Test agent processing
	t.Run("Agent Processing", func(t *testing.T) {
		// Create test server to simulate orchestrator
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" && r.URL.Path == "/internal/task" {
				// Return a test task
				task := map[string]interface{}{
					"id":             1,
					"arg1":           2.0,
					"arg2":           3.0,
					"operation":      "+",
					"operation_time": 100,
				}
				json.NewEncoder(w).Encode(map[string]interface{}{"task": task})
				return
			}

			if r.Method == "POST" && r.URL.Path == "/internal/task" {
				// Verify task result
				var result map[string]interface{}
				json.NewDecoder(r.Body).Decode(&result)
				assert.Equal(t, float64(1), result["id"])
				assert.Equal(t, float64(5), result["result"])
				w.WriteHeader(http.StatusOK)
				return
			}

			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		// Update agent's orchestrator URL
		os.Setenv("ORCHESTRATOR_URL", server.URL)

		// Wait for agent to process task
		time.Sleep(200 * time.Millisecond)
	})

	// Test agent error handling
	t.Run("Agent Error Handling", func(t *testing.T) {
		// Create test server that returns errors
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		// Update agent's orchestrator URL
		os.Setenv("ORCHESTRATOR_URL", server.URL)

		// Wait for agent to handle error
		time.Sleep(100 * time.Millisecond)
	})

	// Test agent parallel processing
	t.Run("Agent Parallel Processing", func(t *testing.T) {
		taskCount := 0
		resultCount := 0

		// Create test server that simulates multiple tasks
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" && r.URL.Path == "/internal/task" {
				taskCount++
				if taskCount <= 2 {
					task := map[string]interface{}{
						"id":             taskCount,
						"arg1":           2.0,
						"arg2":           3.0,
						"operation":      "+",
						"operation_time": 100,
					}
					json.NewEncoder(w).Encode(map[string]interface{}{"task": task})
					return
				}
				w.WriteHeader(http.StatusNotFound)
				return
			}

			if r.Method == "POST" && r.URL.Path == "/internal/task" {
				resultCount++
				w.WriteHeader(http.StatusOK)
				return
			}

			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		// Update agent's orchestrator URL
		os.Setenv("ORCHESTRATOR_URL", server.URL)

		// Wait for agent to process tasks
		time.Sleep(300 * time.Millisecond)

		// Verify that both tasks were processed
		assert.Equal(t, 2, taskCount)
		assert.Equal(t, 2, resultCount)
	})
}
