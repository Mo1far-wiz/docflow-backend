package main

import (
	"docflow-backend/db"
	"docflow-backend/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	db.InitDB()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server := gin.Default()

	routes.RegisterRouts(server)

	server.Run(":8080")
}
