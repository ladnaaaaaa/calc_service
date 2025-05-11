package orchestrator

import (
	"path/filepath"

	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/ladnaaaaaa/calc_service/internal/handlers"
	"github.com/ladnaaaaaa/calc_service/internal/middleware"
)

type Server struct {
	Engine *gin.Engine
	store  *Store
}

func NewServer() *Server {
	engine := gin.Default()

	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(b)))
	templatePath := filepath.Join(projectRoot, "web", "templates", "*")

	engine.LoadHTMLGlob(templatePath)
	engine.Static("/static", filepath.Join(projectRoot, "web", "static"))

	store := NewStore()

	server := &Server{
		Engine: engine,
		store:  store,
	}

	server.setupRoutes()

	return server
}

func (s *Server) Start(addr string) error {
	return s.Engine.Run(addr)
}

func (s *Server) setupRoutes() {
	r := s.Engine

	r.GET("/", s.handleHome)
	r.GET("/expressions", s.handleGetExpressionsRequest)

	r.POST("/api/v1/register", handlers.Register)
	r.POST("/api/v1/login", handlers.Login)

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/calculate", s.handleCalculate)
		protected.GET("/expressions", s.handleGetExpressions)
		protected.GET("/expressions/:id", s.handleGetExpression)
	}

	internal := r.Group("/internal")
	{
		internal.GET("/task", s.handleGetTask)
		internal.POST("/task", s.handleSubmitTask)
	}
}
