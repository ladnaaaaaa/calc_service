package agent

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAgentExecution(t *testing.T) {
	orchestrator := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{
				"task": {
					"id": "test_task",
					"arg1_result": 5,
					"arg2_result": 3,
					"operation": "+",
					"operation_time": 100
				}
			}`))
		}),
	)
	defer orchestrator.Close()

	agent := NewAgent(orchestrator.URL, 2)

	t.Run("execute addition", func(t *testing.T) {
		task := &Task{
			Arg1Result:    5,
			Arg2Result:    3,
			Operation:     "+",
			OperationTime: 100,
		}
		result := agent.executeTask(task)
		assert.Equal(t, 8.0, result)
	})

	t.Run("handle division by zero", func(t *testing.T) {
		task := &Task{
			Arg1Result:    5,
			Arg2Result:    0,
			Operation:     "/",
			OperationTime: 100,
		}
		result := agent.executeTask(task)
		assert.Equal(t, 0.0, result)
	})
}
