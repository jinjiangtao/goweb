package database

import (
	"biaodan/models"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	if err := os.MkdirAll("./data", 0755); err != nil {
		return err
	}

	db, err := gorm.Open(sqlite.Open("./data/biaodan.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db

	return DB.AutoMigrate(
		&models.FormTemplate{},
		&models.FormSubmission{},
	)
}
