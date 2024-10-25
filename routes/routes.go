package routes

import (
	"docflow-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRouts(server *gin.Engine) {
	server.GET("/events", getEvents)
	server.GET("/events/:id", getEvent)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)
	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.DELETE("/events/:id/register", cancelRegistration)

	// POST : generate doc (фак, спеца, рік навчання, тип дока + юзер айді)
	// GET : docs (юзер айді)
	authenticated.POST("/doc/generate", generateDocForUser)
	authenticated.GET("/doc/:id", getDocByID)
	authenticated.GET("/doc/user/:id", getDocsForUser)

	server.POST("/signup", signup)
	server.POST("/login", login)
}
