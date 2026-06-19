package main

import (
	"biaodan/database"
	"biaodan/handlers"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := os.MkdirAll("./uploads", 0755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
	}))

	r.Static("/uploads", "./uploads")
	r.StaticFile("/", "./web/dist/index.html")
	r.Static("/assets", "./web/dist/assets")

	api := r.Group("/api")
	{
		templates := api.Group("/templates")
		{
			templates.GET("", handlers.ListTemplates)
			templates.GET("/:id", handlers.GetTemplate)
			templates.POST("", handlers.CreateTemplate)
			templates.PUT("/:id", handlers.UpdateTemplate)
			templates.DELETE("/:id", handlers.DeleteTemplate)
		}

		submissions := api.Group("/submissions")
		{
			submissions.POST("", handlers.CreateSubmission)
			submissions.GET("", handlers.ListSubmissions)
			submissions.GET("/export", handlers.ExportSubmissions)
			submissions.DELETE("/:id", handlers.DeleteSubmission)
		}

		api.POST("/upload", handlers.UploadFile)

		stats := api.Group("/stats")
		{
			stats.GET("/:templateId", handlers.GetTemplateStats)
		}
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
