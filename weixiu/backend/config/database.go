package config

import (
	"log"
	"weixiu/models"

	sqlite "github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var JWTSecret = []byte("weixiu-oms-2024-secret-key")

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data/weixiu.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Device{},
		&models.WorkOrder{},
		&models.RepairLog{},
		&models.OrderEvaluation{},
		&models.OperationLog{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database initialized successfully")
}
