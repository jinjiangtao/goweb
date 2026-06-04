package main

import (
	"github.com/gin-gonic/gin"
	"student-signup-server/models"
	"student-signup-server/routes"
)

func main() {
	models.InitDB()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	routes.SetupRoutes(r)

	r.Run(":8080")
}
