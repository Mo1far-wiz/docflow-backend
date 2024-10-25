package routes

import (
	"docflow-backend/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRouts(server *gin.Engine) {

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.GET("/user/:id", getUserByID)
	authenticated.GET("/user", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "JWT token is valid."})
	})
	authenticated.POST("/doc/generate", generateDocForUser)
	authenticated.GET("/doc/:id", getDocByID)
	authenticated.GET("/doc/user/:id", getDocsForUser)

	server.POST("/sign-up", signup)
	server.POST("/sign-in", signin)
}
