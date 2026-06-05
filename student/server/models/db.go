package models

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("student.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	err = DB.AutoMigrate(&Signup{}, &Admin{}, &School{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	CreateDefaultAdmin()
}
