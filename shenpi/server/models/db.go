package models

import (
	"log"
	"shenpi/config"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open(config.AppConfig.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&User{}, &ProcessTemplate{}, &ApprovalRequest{}, &ApprovalRecord{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	InitSeedData()

	log.Println("Database initialized successfully")
}

func InitSeedData() {
	var count int64
	DB.Model(&User{}).Count(&count)
	if count > 0 {
		return
	}

	users := []User{
		{Username: "admin", Password: "123456", Name: "系统管理员", Email: "admin@example.com", Role: "admin"},
		{Username: "zhangsan", Password: "123456", Name: "张三", Email: "zhangsan@example.com", Role: "user"},
		{Username: "lisi", Password: "123456", Name: "李四", Email: "lisi@example.com", Role: "user"},
		{Username: "wangwu", Password: "123456", Name: "王五", Email: "wangwu@example.com", Role: "user"},
		{Username: "zhaoliu", Password: "123456", Name: "赵六", Email: "zhaoliu@example.com", Role: "user"},
	}

	for i := range users {
		users[i].HashPassword()
		DB.Create(&users[i])
	}

	log.Println("Seed data initialized successfully")
}
