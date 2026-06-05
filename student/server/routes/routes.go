package routes

import (
	"student-signup-server/handlers"
	"student-signup-server/middleware"

	"github.com/gin-gonic/gin"
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
			protected.POST("/signups", handlers.AdminCreateSignup)
			protected.PUT("/signups/:id", handlers.AdminUpdateSignup)
			protected.PUT("/signups/:id/status", handlers.UpdateSignupStatus)
			protected.GET("/signups/export", handlers.ExportSignups)
			protected.POST("/signups/import", handlers.ImportSignups)
			protected.GET("/stats", handlers.GetStats)

			protected.GET("/schools", handlers.GetSchools)
			protected.GET("/schools/all", handlers.GetAllSchools)
			protected.POST("/schools", handlers.CreateSchool)
			protected.GET("/schools/:id", handlers.GetSchool)
			protected.PUT("/schools/:id", handlers.UpdateSchool)
			protected.DELETE("/schools/:id", handlers.DeleteSchool)
		}
	}
}
