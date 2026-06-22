package routes

import (
	"xiangce/controllers"
	"xiangce/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		share := api.Group("/share")
		{
			share.GET("/:token", controllers.GetShareByToken)
			share.POST("/:token/verify", controllers.VerifySharePassword)
			share.POST("/:token/visit", controllers.RecordVisit)
		}

		files := api.Group("/files")
		{
			files.GET("/:filename", controllers.ServeFile)
			files.GET("/thumb/:filename", controllers.ServeThumbnail)
		}

		authApi := api.Group("")
		authApi.Use(middleware.AuthMiddleware())
		{
			user := authApi.Group("/user")
			{
				user.GET("/info", controllers.GetUserInfo)
				user.GET("/stats", controllers.GetUserStats)
			}

			albums := authApi.Group("/albums")
			{
				albums.GET("", controllers.GetAlbums)
				albums.GET("/:id", controllers.GetAlbum)
				albums.POST("", controllers.CreateAlbum)
				albums.PUT("/:id", controllers.UpdateAlbum)
				albums.DELETE("/:id", controllers.DeleteAlbum)
				albums.POST("/sort", controllers.UpdateAlbumSort)
			}

			media := authApi.Group("/media")
			{
				media.POST("/upload", controllers.UploadMedia)
				media.GET("/album/:albumId", controllers.GetMediaList)
				media.GET("/:id", controllers.GetMedia)
				media.DELETE("/:id", controllers.DeleteMedia)
				media.GET("/:id/download", controllers.DownloadMedia)
				media.POST("/batch/delete", controllers.BatchDeleteMedia)
				media.POST("/batch/download", controllers.BatchDownloadMedia)
				media.POST("/move", controllers.MoveMedia)
				media.POST("/sort/:albumId", controllers.UpdateMediaSort)
			}

			shares := authApi.Group("/shares")
			{
				shares.GET("", controllers.GetShareLinks)
				shares.GET("/album/:albumId", controllers.GetShareLinks)
				shares.POST("", controllers.CreateShareLink)
				shares.DELETE("/:id", controllers.DeleteShareLink)
				shares.GET("/:shareId/visits", controllers.GetVisitRecords)
			}
		}
	}

	return r
}
