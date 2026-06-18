package routes

import (
	"richeng-server/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		events := api.Group("/events")
		{
			events.GET("", handlers.GetEvents)
			events.GET("/:id", handlers.GetEvent)
			events.POST("", handlers.CreateEvent)
			events.PUT("/:id", handlers.UpdateEvent)
			events.DELETE("/:id", handlers.DeleteEvent)
		}

		categories := api.Group("/categories")
		{
			categories.GET("", handlers.GetCategories)
			categories.POST("", handlers.CreateCategory)
			categories.PUT("/:id", handlers.UpdateCategory)
			categories.DELETE("/:id", handlers.DeleteCategory)
		}
	}
}
