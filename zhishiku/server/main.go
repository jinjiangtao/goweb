package main

import (
	"log"

	"zhishiku-server/database"
	"zhishiku-server/handlers"
	"zhishiku-server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	handlers.SetDB(database.DB)
	handlers.InitDefaultAdmin()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		categories := api.Group("/categories")
		{
			categories.GET("", handlers.GetCategories)
		}

		articles := api.Group("/articles")
		{
			articles.GET("", handlers.GetArticles)
			articles.GET("/search", handlers.SearchArticles)
			articles.GET("/:id", handlers.GetArticle)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/auth/me", handlers.GetCurrentUser)

			protected.POST("/articles", handlers.CreateArticle)
			protected.PUT("/articles/:id", handlers.UpdateArticle)
			protected.DELETE("/articles/:id", handlers.DeleteArticle)
			protected.GET("/articles/mine", handlers.GetMyArticles)

			adminRoutes := protected.Group("/admin")
			adminRoutes.Use(middleware.AdminMiddleware())
			{
				adminRoutes.POST("/categories", handlers.CreateCategory)
				adminRoutes.PUT("/categories/:id", handlers.UpdateCategory)
				adminRoutes.DELETE("/categories/:id", handlers.DeleteCategory)
			}
		}
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
