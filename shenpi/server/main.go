package main

import (
	"log"
	"shenpi/config"
	"shenpi/models"
	"shenpi/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()

	models.InitDB()

	r := gin.Default()

	routes.RegisterRoutes(r)

	log.Printf("Server starting on port %s", config.AppConfig.Server.Port)
	if err := r.Run(":" + config.AppConfig.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
