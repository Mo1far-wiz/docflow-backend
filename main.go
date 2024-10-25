package main

import (
	"docflow-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	routes.RegisterRouts(server)

	server.Run(":8080")
}
