package database

import (
	"gongshi/models"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	os.MkdirAll("./data", 0755)
	var err error
	DB, err = gorm.Open(sqlite.Open("./data/gongshi.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.Project{}, &models.WorkRecord{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	seedData()
}

func seedData() {
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		users := []models.User{
			{Username: "admin", Name: "系统管理员", Role: "admin", Password: "123456"},
			{Username: "zhangsan", Name: "张三", Role: "employee", Password: "123456"},
			{Username: "lisi", Name: "李四", Role: "employee", Password: "123456"},
		}
		DB.Create(&users)
	}

	var projectCount int64
	DB.Model(&models.Project{}).Count(&projectCount)
	if projectCount == 0 {
		projects := []models.Project{
			{Name: "ERP系统开发", Code: "PRJ-001", Description: "企业资源规划系统开发项目", Status: "active"},
			{Name: "移动端App开发", Code: "PRJ-002", Description: "公司移动端应用开发", Status: "active"},
			{Name: "数据中台建设", Code: "PRJ-003", Description: "企业数据中台系统建设", Status: "active"},
			{Name: "运维支持", Code: "PRJ-004", Description: "日常系统运维支持工作", Status: "active"},
		}
		DB.Create(&projects)
	}
}
