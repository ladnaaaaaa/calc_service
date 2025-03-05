package orchestrator

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	store  *Store
}

func NewServer() *Server {
	return &Server{
		engine: gin.Default(),
		store:  NewStore(),
	}
}

func (s *Server) Start(addr string) error {
	s.engine.LoadHTMLGlob("web/templates/*")
	s.setupRoutes()
	s.registerWebRoutes()
	return s.engine.Run(addr)
}

func (s *Server) setupRoutes() {
	s.engine.POST("/api/v1/calculate", s.handleCalculate)
	s.engine.GET("/api/v1/expressions", s.handleGetExpressions)
	s.engine.GET("/api/v1/expressions/:id", s.handleGetExpression)
	s.engine.GET("/internal/task", s.handleGetTask)
	s.engine.POST("/internal/task", s.handleSubmitTask)
}
