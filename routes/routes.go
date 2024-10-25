package routes

import (
	"docflow-backend/middlewares"
	"docflow-backend/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRouts(server *gin.Engine) {

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.GET("/user/:id", getUserByID)
	authenticated.GET("/user", func(context *gin.Context) {
		userId := context.GetInt64("userId")
		user, err := models.GeUserByID(userId)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get user by ID."})
		}
		user.Password = ""
		fmt.Println(user)
		context.JSON(http.StatusOK, gin.H{"message": "JWT token is valid.", "user": user})
	})
	authenticated.POST("/doc/generate", generateDocForUser)
	authenticated.GET("/doc/:id", getDocByID)
	authenticated.GET("/doc/user/:id", getDocsForUser)

	server.POST("/sign-up", signup)
	server.POST("/sign-in", signin)
}
