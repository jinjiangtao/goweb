package routes

import (
	"biji/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/stats", handlers.GetStats)

		notes := api.Group("/notes")
		{
			notes.GET("", handlers.GetNotes)
			notes.POST("", handlers.CreateNote)
			notes.GET("/:id", handlers.GetNote)
			notes.PUT("/:id", handlers.UpdateNote)
			notes.DELETE("/:id", handlers.DeleteNote)
			notes.POST("/:id/verify", handlers.VerifyNotePassword)
			notes.POST("/:id/restore", handlers.RestoreNote)
			notes.POST("/:id/pin", handlers.TogglePin)
			notes.POST("/:id/archive", handlers.ToggleArchive)
		}

		api.DELETE("/trash/empty", handlers.EmptyTrash)

		tags := api.Group("/tags")
		{
			tags.GET("", handlers.GetTags)
			tags.POST("", handlers.CreateTag)
			tags.GET("/:id", handlers.GetTag)
			tags.PUT("/:id", handlers.UpdateTag)
			tags.DELETE("/:id", handlers.DeleteTag)
			tags.GET("/:id/count", handlers.GetTagNoteCount)
			tags.POST("/reorder", handlers.ReorderTags)
		}

		groups := api.Group("/tag-groups")
		{
			groups.GET("", handlers.GetTagGroups)
			groups.POST("", handlers.CreateTagGroup)
			groups.PUT("/:id", handlers.UpdateTagGroup)
			groups.DELETE("/:id", handlers.DeleteTagGroup)
		}
	}
}
