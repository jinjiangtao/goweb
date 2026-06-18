package database

import (
	"log"
	"richeng-server/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("richeng.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	err = DB.AutoMigrate(&models.Category{}, &models.Event{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	seedData()
}

func seedData() {
	var count int64
	DB.Model(&models.Category{}).Count(&count)
	if count == 0 {
		categories := []models.Category{
			{Name: "工作", Color: "#3B82F6"},
			{Name: "生活", Color: "#10B981"},
			{Name: "学习", Color: "#F59E0B"},
			{Name: "运动", Color: "#EF4444"},
			{Name: "其他", Color: "#8B5CF6"},
		}
		DB.Create(&categories)
		log.Println("Seed categories created")
	}
}
