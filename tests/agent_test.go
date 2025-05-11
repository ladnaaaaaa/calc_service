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
	os.Setenv("ORCHESTRATOR_URL", "http://localhost:8080")
	os.Setenv("COMPUTING_POWER", "2")
	os.Setenv("TIME_ADDITION_MS", "100")
	os.Setenv("TIME_SUBTRACTION_MS", "100")
	os.Setenv("TIME_MULTIPLICATIONS_MS", "100")
	os.Setenv("TIME_DIVISIONS_MS", "100")

	return agent.NewAgent("http://localhost:8080", 2)
}

func TestAgent(t *testing.T) {
	agent := setupTestAgent()

	go agent.Start()

	time.Sleep(100 * time.Millisecond)

	t.Run("Agent Processing", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" && r.URL.Path == "/internal/task" {
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

		os.Setenv("ORCHESTRATOR_URL", server.URL)

		time.Sleep(200 * time.Millisecond)
	})

	t.Run("Agent Error Handling", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		os.Setenv("ORCHESTRATOR_URL", server.URL)

		time.Sleep(100 * time.Millisecond)
	})

	t.Run("Agent Parallel Processing", func(t *testing.T) {
		taskCount := 0
		resultCount := 0

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

		os.Setenv("ORCHESTRATOR_URL", server.URL)

		time.Sleep(300 * time.Millisecond)

		assert.Equal(t, 2, taskCount)
		assert.Equal(t, 2, resultCount)
	})
}
