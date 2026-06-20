package main

import (
	"logistics-server/database"
	"logistics-server/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.POST("/orders", handlers.CreateOrder)
		api.GET("/orders", handlers.GetOrders)
		api.GET("/orders/search", handlers.SearchOrder)
		api.GET("/orders/:id", handlers.GetOrder)
		api.PUT("/orders/:id", handlers.UpdateOrder)
		api.DELETE("/orders/:id", handlers.DeleteOrder)

		api.POST("/orders/:id/nodes", handlers.AddNode)
		api.GET("/orders/:id/nodes", handlers.GetNodes)
		api.GET("/orders/:id/progress", handlers.GetProgress)

		api.PUT("/nodes/:node_id", handlers.UpdateNode)
		api.DELETE("/nodes/:node_id", handlers.DeleteNode)
	}

	r.Run(":8080")
}
