package main

import (
	"log"
	"neitui/models"
	"neitui/routes"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	uploadsPath := "./uploads/resumes"
	if _, err := os.Stat(uploadsPath); os.IsNotExist(err) {
		os.MkdirAll(uploadsPath, 0755)
	}

	r.Static("/uploads", "./uploads")

	if err := models.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
