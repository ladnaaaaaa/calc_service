package orchestrator

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ladnaaaaaa/calc_service/internal/database"
	"github.com/ladnaaaaaa/calc_service/internal/models"
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

func (s *Store) IsTaskReady(task *models.Task) bool {
	// Get all tasks for this expression
	tasks, err := s.GetTasksByExpressionID(task.ExpressionID)
	if err != nil {
		return false
	}

	// Check if all tasks with lower order numbers are completed
	for _, t := range tasks {
		if t.OrderNum < task.OrderNum && t.Status != models.StatusCompleted {
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

func (s *Store) AddExpression(expr *models.Expression) error {
	return database.DB.Create(expr).Error
}

func (s *Store) GetExpression(id uint) (*models.Expression, error) {
	var expr models.Expression
	err := database.DB.Preload("Tasks").First(&expr, id).Error
	if err != nil {
		return nil, err
	}
	return &expr, nil
}

func (s *Store) GetAllExpressions(userID uint) ([]models.Expression, error) {
	var expressions []models.Expression
	err := database.DB.Where("user_id = ?", userID).Preload("Tasks").Find(&expressions).Error
	return expressions, err
}

func (s *Store) AddTask(task *models.Task) error {
	return database.DB.Create(task).Error
}

func (s *Store) GetTask(id uint) (*models.Task, error) {
	var task models.Task
	err := database.DB.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *Store) UpdateTask(task *models.Task) error {
	return database.DB.Save(task).Error
}

func (s *Store) GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task
	err := database.DB.Find(&tasks).Error
	return tasks, err
}

func (s *Store) GetTasksByExpressionID(exprID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := database.DB.Where("expression_id = ?", exprID).Find(&tasks).Error
	return tasks, err
}

func (s *Store) GetPendingTasks() ([]models.Task, error) {
	var tasks []models.Task
	err := database.DB.Where("status = ?", models.StatusPending).Find(&tasks).Error
	return tasks, err
}

func (s *Store) RemoveTask(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tasks, id)
}

func (s *Store) GetOperationTime(op string) time.Duration {
	return s.opTimes[op]
}

func (s *Store) UpdateExpression(expr *models.Expression) error {
	return database.DB.Save(expr).Error
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
	OrderNum     int
}
