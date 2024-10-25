package routes

import (
	"net/http"
	"strconv"

	"docflow-backend/models"
	"docflow-backend/utils"

	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully."})
}

func signin(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = user.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authorize user."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Logged successfully.", "token": token})
}

func getUserByID(context *gin.Context) {
	userId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse user ID."})
		return
	}

	user, err := models.GeUserByID(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get user by ID."})
		return
	}
	user.Password = ""

	context.JSON(http.StatusOK, user)
}
