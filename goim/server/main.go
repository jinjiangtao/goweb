package main

import (
	"goim/server/cache"
	"goim/server/handler"
	"goim/server/storage"
	"goim/server/websocket"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	storage.InitDB()
	cache.InitCache()

	r := gin.Default()

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
		}

		user := api.Group("/user")
		{
			user.POST("/add-friend", handler.AddFriend)
			user.GET("/friends", handler.GetFriends)
			user.GET("/online-status", handler.GetOnlineStatus)
			user.GET("/online-users", handler.GetOnlineUsers)
			user.GET("/profile", handler.GetUserProfile)
			user.POST("/upload-avatar", handler.UploadAvatar)
			user.POST("/update-profile", handler.UpdateProfile)
		}

		group := api.Group("/group")
		{
			group.POST("/create", handler.CreateGroup)
			group.POST("/add-member", handler.AddGroupMember)
			group.POST("/remove-member", handler.RemoveGroupMember)
			group.GET("/list", handler.GetGroups)
			group.GET("/members/:groupID", handler.GetGroupMembers)
		}

		message := api.Group("/message")
		{
			message.GET("/history", handler.GetMessageHistory)
			message.POST("/read", handler.MarkMessageRead)
		}

		// Admin API
		admin := api.Group("/admin")
		{
			admin.POST("/login", handler.AdminLogin)
			admin.Use(handler.AdminAuthMiddleware())
			admin.GET("/users", handler.GetAllUsers)
			admin.GET("/messages", handler.GetAllMessages)
		}
	}

	r.GET("/ws/:userID", func(c *gin.Context) {
		userID := c.Param("userID")
		c.Request.URL.RawQuery = "userID=" + userID
		websocket.HandleWebSocket(c.Writer, c.Request)
	})

	go websocket.GetHub().Run()

	log.Println("Server started on :8080")
	log.Fatal(r.Run(":8080"))
}
