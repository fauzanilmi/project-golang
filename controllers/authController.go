package controllers

import (
	"fmt"
	"net/http"
	"project/models"
	"project/token"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

func Register(ctx *gin.Context) {
	var input RegisterInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c := models.Customer{}

	c.Username = input.Username
	c.Password = input.Password
	c.Name = input.Name

	_, err := c.SaveUser()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Registration success\n Hello, " + input.Name})
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(ctx *gin.Context) {
	var input LoginInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c := models.Customer{}

	c.Username = input.Username
	c.Password = input.Password

	token, customer, err := models.LoginCheck(c.Username, c.Password)

	if err != nil {
		if token != "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": token})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "Hello, " + customer.Name + "!"})

}

func Logout(ctx *gin.Context) {
	auth, err := token.ExtractTokenAuth(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	name, delleteErr := models.DeleteAuth(auth)
	if delleteErr != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	message := fmt.Sprintf("Goodbye, %s!", name)

	ctx.JSON(http.StatusOK, gin.H{"message": message})
}
