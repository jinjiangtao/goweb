package server

import (
	"log"
	"os"

	"whiteboard/internal/handlers"
	"whiteboard/internal/middleware"
	"whiteboard/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Server struct {
	engine *gin.Engine
	port   string
}

func New() *Server {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	models.EnsureDataDir()
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./data/whiteboard.db"
	}
	err = models.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	engine.Use(middleware.Logger())
	engine.Use(middleware.Recovery())
	engine.Use(middleware.CORS())

	server := &Server{
		engine: engine,
		port:   os.Getenv("SERVER_PORT"),
	}

	if server.port == "" {
		server.port = "8080"
	}

	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	api := s.engine.Group("/api")
	{
		boards := api.Group("/boards")
		{
			boards.GET("", handlers.ListBoards)
			boards.POST("", handlers.CreateBoard)
			boards.GET("/:id", handlers.GetBoard)
			boards.PUT("/:id", handlers.UpdateBoard)
			boards.DELETE("/:id", handlers.DeleteBoard)
			boards.GET("/:id/elements", handlers.GetElements)
			boards.POST("/:id/elements", handlers.SaveElements)
			boards.POST("/:id/clear", handlers.ClearBoard)
			boards.GET("/:id/history", handlers.GetHistory)
			boards.POST("/:id/history/:historyId/restore", handlers.RestoreHistory)
		}
	}

	s.engine.GET("/ws/:id", handlers.HandleWebSocket)

	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Whiteboard server is running",
		})
	})
}

func (s *Server) Run() error {
	log.Printf("Server starting on port %s...", s.port)
	log.Printf("API: http://localhost:%s/api", s.port)
	log.Printf("WS:  ws://localhost:%s/ws/:id", s.port)
	return s.engine.Run(":" + s.port)
}
