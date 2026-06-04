package routes

import (
	"github.com/gin-gonic/gin"
	"student-signup-server/handlers"
	"student-signup-server/middleware"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/api/signup", handlers.Signup)

	admin := r.Group("/api/admin")
	{
		admin.POST("/login", handlers.Login)

		protected := admin.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/signups", handlers.GetSignups)
			protected.PUT("/signups/:id/status", handlers.UpdateSignupStatus)
			protected.GET("/stats", handlers.GetStats)
		}
	}
}
