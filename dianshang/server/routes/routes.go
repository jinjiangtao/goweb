package routes

import (
	"server/handlers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())

	api := r.Group("/api/admin")
	{
		api.POST("/login", handlers.Login)
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/info", handlers.GetInfo)
			auth.POST("/logout", handlers.Logout)

			users := auth.Group("/users")
			users.GET("", handlers.GetUsers)
			users.POST("", middleware.SuperAdminMiddleware(), handlers.CreateUser)
			users.PUT("/:id", middleware.AdminMiddleware(), handlers.UpdateUser)
			users.DELETE("/:id", middleware.SuperAdminMiddleware(), handlers.DeleteUser)
			users.PUT("/:id/status", middleware.AdminMiddleware(), handlers.UpdateUserStatus)
			users.PUT("/:id/password", middleware.AdminMiddleware(), handlers.ResetPassword)

			menus := auth.Group("/menus")
			menus.GET("", handlers.GetMenus)
			menus.POST("", middleware.AdminMiddleware(), handlers.CreateMenu)
			menus.PUT("/:id", middleware.AdminMiddleware(), handlers.UpdateMenu)
			menus.DELETE("/:id", middleware.AdminMiddleware(), handlers.DeleteMenu)

			roles := auth.Group("/roles")
			roles.GET("", handlers.GetRoles)
			roles.GET("/:role/menus", handlers.GetRoleMenus)
			roles.PUT("/:role/menus", middleware.SuperAdminMiddleware(), handlers.SetRoleMenus)
		}
	}
	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}