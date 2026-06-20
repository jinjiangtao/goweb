package main

import (
	"jianli-server/config"
	"jianli-server/db"
	"jianli-server/middleware"
	"jianli-server/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	db.InitDB()

	r := gin.Default()

	r.Use(middleware.CORS())

	routes.SetupRoutes(r)

	log.Printf("Server starting on port %s...", config.AppConfig.ServerPort)
	if err := r.Run(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
