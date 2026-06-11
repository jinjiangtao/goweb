package main

import (
	"erp/database"
	"erp/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	r := gin.Default()

	routes.SetupRoutes(r)

	fmt.Println("Server starting on port 8080...")
	r.Run(":8080")
}
