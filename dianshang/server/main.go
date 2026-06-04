package main

import (
	"server/models"
	"server/routes"
)

func main() {
	models.InitDB()
	models.InitSuperAdmin()
	r := routes.SetupRouter()
	r.Run(":8081")
}
