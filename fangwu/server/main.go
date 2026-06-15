package main

import (
	"log"

	"fangwu-server/config"
	"fangwu-server/handlers"
	"fangwu-server/middleware"
	"fangwu-server/models"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	config.Init()

	db, err := gorm.Open(sqlite.Open(config.Cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败: ", err)
	}
	models.DB = db

	db.AutoMigrate(&models.Admin{}, &models.House{})

	seedAdmin(db)

	r := gin.Default()

	r.Use(middleware.CORS())

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		api.POST("/login", handlers.Login)
	}

	authAPI := api.Group("")
	authAPI.Use(middleware.AuthRequired())
	{
		authAPI.GET("/profile", handlers.GetProfile)
		authAPI.PUT("/change-password", handlers.ChangePassword)

		authAPI.GET("/houses", handlers.ListHouses)
		authAPI.GET("/houses/:id", handlers.GetHouse)
		authAPI.POST("/houses", handlers.CreateHouse)
		authAPI.PUT("/houses/:id", handlers.UpdateHouse)
		authAPI.DELETE("/houses/:id", handlers.DeleteHouse)
		authAPI.PUT("/houses/:id/status", handlers.ToggleStatus)
		authAPI.PUT("/houses/:id/recommend", handlers.ToggleRecommend)

		authAPI.POST("/upload", handlers.UploadImage)
	}

	log.Printf("服务启动于端口 %s", config.Cfg.Port)
	if err := r.Run(":" + config.Cfg.Port); err != nil {
		log.Fatal("服务启动失败: ", err)
	}
}

func seedAdmin(db *gorm.DB) {
	var count int64
	db.Model(&models.Admin{}).Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		admin := models.Admin{
			Username: "admin",
			Password: string(hashedPassword),
		}
		db.Create(&admin)
		log.Println("默认管理员账号已创建: admin / 123456")
	}
}
