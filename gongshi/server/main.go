package main

import (
	"gongshi/database"
	"gongshi/handlers"
	"log"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetHeader("X-User-Role")
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.Request.Header.Set("X-User-ID", "2")
			c.Request.Header.Set("X-User-Role", "employee")
		}
		if userRole == "" {
			c.Request.Header.Set("X-User-Role", "employee")
		}
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("X-User-Role")
		if !strings.EqualFold(role, "admin") {
			c.JSON(403, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	database.InitDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	api.Use(AuthMiddleware())
	{
		api.POST("/login", handlers.Login)
		api.GET("/user/current", handlers.GetCurrentUser)
		api.GET("/users", handlers.GetUsers)

		projects := api.Group("/projects")
		{
			projects.GET("", handlers.GetProjects)
			projects.POST("", AdminMiddleware(), handlers.CreateProject)
			projects.PUT("/:id", AdminMiddleware(), handlers.UpdateProject)
			projects.DELETE("/:id", AdminMiddleware(), handlers.DeleteProject)
		}

		records := api.Group("/work-records")
		{
			records.GET("", handlers.GetWorkRecords)
			records.POST("", handlers.CreateWorkRecord)
			records.POST("/batch", handlers.BatchCreateWorkRecords)
			records.PUT("/:id", handlers.UpdateWorkRecord)
			records.DELETE("/:id", handlers.DeleteWorkRecord)
			records.POST("/:id/approve", AdminMiddleware(), handlers.ApproveWorkRecord)
			records.POST("/:id/reject", AdminMiddleware(), handlers.RejectWorkRecord)
			records.POST("/batch-approve", AdminMiddleware(), handlers.BatchApprove)
		}

		stats := api.Group("/stats")
		{
			stats.GET("/daily", handlers.GetDailyHoursSummary)
			stats.GET("/project", handlers.GetProjectHoursSummary)
			stats.GET("/overall", handlers.GetOverallStats)
		}
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
