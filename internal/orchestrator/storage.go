package orchestrator

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type Store struct {
	expressions map[string]*Expression
	tasks       map[string]*Task
	mu          sync.RWMutex
	opTimes     map[string]time.Duration
}

func NewStore() *Store {
	store := &Store{
		expressions: make(map[string]*Expression),
		tasks:       make(map[string]*Task),
		opTimes:     make(map[string]time.Duration),
	}

	store.opTimes["+"] = parseDuration("TIME_ADDITION_MS", 1000)
	store.opTimes["-"] = parseDuration("TIME_SUBTRACTION_MS", 1000)
	store.opTimes["*"] = parseDuration("TIME_MULTIPLICATIONS_MS", 2000)
	store.opTimes["/"] = parseDuration("TIME_DIVISIONS_MS", 3000)

	return store
}

func (s *Store) IsTaskReady(task *Task) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, depID := range task.DependsOn {
		depTask, exists := s.tasks[depID]
		if !exists || depTask.Status != "completed" {
			return false
		}
	}
	return true
}

func parseDuration(envVar string, defaultMs int) time.Duration {
	val := os.Getenv(envVar)
	if val == "" {
		return time.Duration(defaultMs) * time.Millisecond
	}

	ms, err := strconv.Atoi(val)
	if err != nil || ms < 0 {
		return time.Duration(defaultMs) * time.Millisecond
	}
	return time.Duration(ms) * time.Millisecond
}

func (s *Store) AddExpression(expr *Expression) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.expressions[expr.ID] = expr

	for _, task := range expr.Tasks {
		s.tasks[task.ID] = task
	}
}

func (s *Store) GetExpression(id string) (*Expression, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	expr, exists := s.expressions[id]
	return expr, exists
}

func (s *Store) GetAllExpressions() []*Expression {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Expression, 0, len(s.expressions))
	for _, expr := range s.expressions {
		result = append(result, expr)
	}
	return result
}

func (s *Store) AddTask(task *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
}

func (s *Store) GetTask(id string) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[id]
	return task, exists
}

func (s *Store) UpdateTask(task *Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[task.ID]; !exists {
		return fmt.Errorf("task not found")
	}

	s.tasks[task.ID] = task
	return nil
}

func (s *Store) GetAllTasks() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}
	return result
}

func (s *Store) GetTasksByExpressionID(exprID string) []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Task, 0)
	for _, task := range s.tasks {
		if task.ExpressionID == exprID {
			result = append(result, task)
		}
	}
	return result
}

func (s *Store) GetPendingTasks() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Task, 0)
	for _, task := range s.tasks {
		if task.Status == "pending" {
			result = append(result, task)
		}
	}
	return result
}

func (s *Store) RemoveTask(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tasks, id)
}

func (s *Store) GetOperationTime(op string) time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.opTimes[op]
}

type Expression struct {
	ID          string
	Status      string
	Result      float64
	Tasks       []*Task
	FinalTaskID string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Task struct {
	ID           string
	ExpressionID string
	Arg1ID       string // ID первого аргумента
	Arg2ID       string // ID второго аргумента
	Operation    string
	Status       string
	Result       float64
	DependsOn    []string // Массив ID зависимых задач
	StartedAt    time.Time
	CompletedAt  time.Time
}
