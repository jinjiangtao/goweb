package main

import (
	"codesnippet/database"
	"codesnippet/handlers"
	"codesnippet/middleware"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.StaticFS("/uploads", gin.Dir("uploads", true))

	api := r.Group("/api")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
		api.GET("/languages", handlers.GetLanguages)

		api.GET("/snippets", middleware.OptionalAuthMiddleware(), handlers.ListSnippets)
		api.GET("/snippets/search", middleware.OptionalAuthMiddleware(), handlers.SearchSnippets)
		api.GET("/snippets/:id", middleware.OptionalAuthMiddleware(), handlers.GetSnippet)
		api.GET("/snippets/:id/comments", middleware.OptionalAuthMiddleware(), handlers.GetComments)

		api.POST("/diff", handlers.DiffSnippets)

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/user/me", handlers.GetCurrentUser)

			auth.POST("/snippets", handlers.CreateSnippet)
			auth.PUT("/snippets/:id", handlers.UpdateSnippet)
			auth.DELETE("/snippets/:id", handlers.DeleteSnippet)
			auth.GET("/my/snippets", handlers.GetMySnippets)
			auth.GET("/my/favorites", handlers.GetMyFavorites)

			auth.POST("/snippets/:id/like", handlers.ToggleLike)
			auth.POST("/snippets/:id/favorite", handlers.ToggleFavorite)
			auth.POST("/snippets/:id/comments", handlers.CreateComment)
			auth.DELETE("/comments/:comment_id", handlers.DeleteComment)

			auth.POST("/snippets/:id/fork", handlers.ForkSnippet)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
