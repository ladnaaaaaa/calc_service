package agent

import (
	"log"
	"sync"
	"time"
)

type Agent struct {
	client  *Client
	workers int
}

func NewAgent(serverURL string, workers int) *Agent {
	return &Agent{
		client:  NewClient(serverURL),
		workers: workers,
	}
}

func (a *Agent) Start() {
	var wg sync.WaitGroup
	for i := 0; i < a.workers; i++ {
		wg.Add(1)
		go a.worker(&wg)
	}
	wg.Wait()
}

func (a *Agent) worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task, err := a.client.FetchTask()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		result := a.executeTask(task)

		if err := a.client.SubmitResult(task.ID, result); err != nil {
			log.Printf("Failed to submit result: %v", err)
		}
	}
}

func (a *Agent) executeTask(task *Task) float64 {
	operationTime := time.Duration(task.OperationTime) * time.Millisecond
	time.Sleep(operationTime)

	switch task.Operation {
	case "+":
		return task.Arg1Result + task.Arg2Result
	case "-":
		return task.Arg1Result - task.Arg2Result
	case "*":
		return task.Arg1Result * task.Arg2Result
	case "/":
		if task.Arg2Result == 0 {
			return 0
		}
		return task.Arg1Result / task.Arg2Result
	default:
		return 0
	}
}
