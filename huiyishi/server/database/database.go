package database

import (
	"huiyishi-server/models"
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("huiyishi.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB.AutoMigrate(&models.Admin{}, &models.Room{}, &models.Booking{}, &models.User{})

	var count int64
	DB.Model(&models.Admin{}).Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		admin := models.Admin{
			Username:  "admin",
			Password:  string(hashedPassword),
			Nickname:  "系统管理员",
			CreatedAt: time.Now(),
		}
		DB.Create(&admin)
		log.Println("Default admin created: admin/123456")
	}

	DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		user := models.User{
			Username:  "test",
			Password:  string(hashedPassword),
			RealName:  "测试用户",
			Phone:     "13800138000",
			Status:    1,
			CreatedAt: time.Now(),
		}
		DB.Create(&user)
		log.Println("Default test user created: test/123456")
	}
}
