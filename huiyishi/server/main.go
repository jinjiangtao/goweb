
package main

import (
	"log"
	"huiyishi-server/database"
	"huiyishi-server/handlers"
	"huiyishi-server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.POST("/api/admin/login", handlers.Login)

	admin := r.Group("/api/admin")
	admin.Use(middleware.JWTAuth())
	{
		admin.GET("/rooms", handlers.GetRooms)
		admin.POST("/rooms", handlers.CreateRoom)
		admin.PUT("/rooms/:id", handlers.UpdateRoom)
		admin.DELETE("/rooms/:id", handlers.DeleteRoom)

		admin.GET("/bookings", handlers.GetBookings)
		admin.PUT("/bookings/:id/cancel", handlers.CancelBooking)

		admin.GET("/stats", handlers.GetStats)
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
