package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"weixiu/config"
	"weixiu/handlers"
	"weixiu/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}
	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatal("Failed to create uploads directory:", err)
	}

	config.InitDB()
	handlers.InitDefaultUsers()

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			handlers.CheckOverdueOrders()
		}
	}()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./uploads")
	r.Static("/assets", "../frontend/dist/assets")
	r.GET("/", func(c *gin.Context) {
		c.File("../frontend/dist/index.html")
	})
	r.NoRoute(func(c *gin.Context) {
		c.File("../frontend/dist/index.html")
	})

	api := r.Group("/api")
	{
		api.POST("/auth/login", handlers.Login)
		api.POST("/auth/register", handlers.Register)

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/auth/me", handlers.GetCurrentUser)
			auth.GET("/users", middleware.RoleMiddleware("admin"), handlers.GetUsers)
			auth.GET("/technicians", handlers.GetTechnicians)

			orders := auth.Group("/orders")
			{
				orders.POST("", handlers.CreateWorkOrder)
				orders.GET("", handlers.GetWorkOrders)
				orders.GET("/:id", handlers.GetWorkOrder)
				orders.POST("/:id/assign", middleware.RoleMiddleware("admin"), handlers.AssignWorkOrder)
				orders.POST("/:id/auto-assign", middleware.RoleMiddleware("admin"), handlers.AutoAssignWorkOrder)
				orders.POST("/:id/start", middleware.RoleMiddleware("technician", "admin"), handlers.StartProcessOrder)
				orders.POST("/:id/process", middleware.RoleMiddleware("technician", "admin"), handlers.ProcessWorkOrder)
				orders.POST("/:id/complete", middleware.RoleMiddleware("technician", "admin"), handlers.CompleteWorkOrder)
				orders.POST("/:id/reject", middleware.RoleMiddleware("admin"), handlers.RejectWorkOrder)
				orders.POST("/:id/cancel", handlers.CancelWorkOrder)
				orders.POST("/:id/evaluate", handlers.EvaluateWorkOrder)
			}

			devices := auth.Group("/devices")
			devices.Use(middleware.RoleMiddleware("admin", "technician"))
			{
				devices.POST("", handlers.CreateDevice)
				devices.GET("", handlers.GetDevices)
				devices.GET("/:id", handlers.GetDevice)
				devices.PUT("/:id", handlers.UpdateDevice)
				devices.DELETE("/:id", handlers.DeleteDevice)
			}

			stats := auth.Group("/stats")
			stats.Use(middleware.RoleMiddleware("admin"))
			{
				stats.GET("", handlers.GetStatistics)
			}

			logs := auth.Group("/logs")
			logs.Use(middleware.RoleMiddleware("admin"))
			{
				logs.GET("", handlers.GetOperationLogs)
			}

			auth.POST("/upload", func(c *gin.Context) {
				file, err := c.FormFile("file")
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
					return
				}

				ext := filepath.Ext(file.Filename)
				filename := time.Now().Format("20060102150405") + "_" + file.Filename
				savePath := filepath.Join("uploads", filename)

				if err := c.SaveUploadedFile(file, savePath); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"url":  "/uploads/" + filename,
					"name": filename,
					"ext":  ext,
				})
			})
		}
	}

	log.Println("Server starting on :8080")
	log.Println("Default accounts:")
	log.Println("  Admin: admin / 123456")
	log.Println("  User:  user1 / 123456")
	log.Println("  Tech:  tech1 / 123456, tech2 / 123456")
	r.Run(":8080")
}
