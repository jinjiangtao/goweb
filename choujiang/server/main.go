package main

import (
	"choujiang/handlers"
	"choujiang/middleware"
	"choujiang/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("choujiang.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	err = db.AutoMigrate(&models.Admin{}, &models.Prize{}, &models.Record{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	var adminCount int64
	db.Model(&models.Admin{}).Count(&adminCount)
	if adminCount == 0 {
		db.Create(&models.Admin{
			Username: "admin",
			Password: "123456",
		})
	}

	models.InitDB(db)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.POST("/api/login", handlers.Login)
	r.GET("/api/prizes", handlers.GetPublicPrizes)
	r.POST("/api/lottery/draw", handlers.DoLottery)
	r.GET("/api/lottery/records", handlers.GetRecordsByPhone)
	r.PUT("/api/lottery/records/:id/claim", handlers.ClaimRecordPublic)

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/admin/prizes", handlers.GetPrizes)
		auth.POST("/admin/prizes", handlers.CreatePrize)
		auth.PUT("/admin/prizes/:id", handlers.UpdatePrize)
		auth.DELETE("/admin/prizes/:id", handlers.DeletePrize)
		auth.PUT("/admin/prizes/:id/toggle", handlers.TogglePrize)

		auth.GET("/admin/records", handlers.GetRecords)
		auth.PUT("/admin/records/:id/claim", handlers.ClaimRecord)
		auth.GET("/admin/records/export", handlers.ExportRecords)

		auth.GET("/admin/stats", handlers.GetStats)
	}

	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}
