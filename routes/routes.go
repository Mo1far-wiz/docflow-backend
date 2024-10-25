package routes

import (
	"docflow-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRouts(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/doc/generate", generateDocForUser)
	authenticated.GET("/doc/:id", getDocByID)
	authenticated.GET("/doc/user/:id", getDocsForUser)

	server.POST("/sign-up", signup)
	server.POST("/sign-in", signin)
}
