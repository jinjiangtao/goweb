package database

import (
	"codesnippet/models"
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() error {
	dbPath := filepath.Join(".", "data")
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return err
	}

	dbFile := filepath.Join(dbPath, "codesnippet.db")
	var err error
	DB, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Snippet{},
		&models.Like{},
		&models.Favorite{},
		&models.Comment{},
	)
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}
