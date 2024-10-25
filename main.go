package main

import (
	"docflow-backend/db"
	"docflow-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()

	routes.RegisterRouts(server)

	server.Run(":8080")
}
