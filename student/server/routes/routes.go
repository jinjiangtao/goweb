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
			protected.GET("/stats/daily", handlers.GetDailyStats)

			protected.GET("/schools", handlers.GetSchools)
			protected.GET("/schools/all", handlers.GetAllSchools)
			protected.POST("/schools", handlers.CreateSchool)
			protected.GET("/schools/:id", handlers.GetSchool)
			protected.PUT("/schools/:id", handlers.UpdateSchool)
			protected.DELETE("/schools/:id", handlers.DeleteSchool)

			protected.GET("/users", handlers.GetUsers)
			protected.POST("/users", handlers.CreateUser)
			protected.PUT("/users/:id", handlers.UpdateUser)
			protected.DELETE("/users/:id", handlers.DeleteUser)
			protected.POST("/users/:id/reset-password", handlers.ResetPassword)
			protected.GET("/user/info", handlers.GetUserInfo)

			protected.GET("/menus", handlers.GetMenus)
			protected.GET("/menus/tree", handlers.GetMenuTree)
			protected.GET("/menus/parent", handlers.GetParentMenus)
			protected.POST("/menus", handlers.CreateMenu)
			protected.PUT("/menus/:id", handlers.UpdateMenu)
			protected.DELETE("/menus/:id", handlers.DeleteMenu)
		}
	}
}
