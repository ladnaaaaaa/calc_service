package orchestrator

import (
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
	server := &Server{
		Engine: engine,
		store:  NewStore(),
	}

	// Load templates
	engine.LoadHTMLGlob("web/templates/*")

	// Setup all routes
	server.setupRoutes()

	return server
}

func (s *Server) Start(addr string) error {
	return s.Engine.Run(addr)
}

func (s *Server) setupRoutes() {
	r := s.Engine

	// Static files
	r.Static("/static", "web/static")

	// Web routes
	r.GET("/", s.handleHome)
	r.GET("/expressions", s.handleGetExpressionsRequest)

	// Public API routes
	r.POST("/api/v1/register", handlers.Register)
	r.POST("/api/v1/login", handlers.Login)

	// Protected API routes
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/calculate", s.handleCalculate)
		protected.GET("/expressions", s.handleGetExpressions)
		protected.GET("/expressions/:id", s.handleGetExpression)
	}

	// Internal routes for agent communication
	internal := r.Group("/internal")
	{
		internal.GET("/task", s.handleGetTask)
		internal.POST("/task", s.handleSubmitTask)
	}
}
