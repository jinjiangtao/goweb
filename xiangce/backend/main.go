package main

import (
	"log"

	"xiangce/config"
	"xiangce/models"
	"xiangce/routes"
	"xiangce/utils"
)

func main() {
	cfg := config.LoadConfig()

	utils.SetJWTSecret(cfg.JWTSecret)

	err := models.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := routes.SetupRouter()

	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Upload directory: %s", cfg.UploadDir)
	log.Printf("Database: %s", cfg.DBPath)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
