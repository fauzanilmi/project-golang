package main

import (
	"project/controllers"
	"project/middleware"
	"project/models"

	"github.com/gin-gonic/gin"
)

func main() {

	models.ConnectDataBase()
	restAPI := gin.Default()

	authGroup := restAPI.Group("/api/auth")

	authGroup.POST("/register", controllers.Register)
	authGroup.POST("/login", controllers.Login)
	authGroup.POST("/logout", controllers.Logout).Use(middleware.JwtAuthMiddleware())

	customerGroup := restAPI.Group("/api/customer")
	customerGroup.Use(middleware.JwtAuthMiddleware())
	customerGroup.GET("/balance", controllers.CheckBalance)
	customerGroup.POST("/topup", controllers.TopUpBalance)
	customerGroup.POST("/withdraw", controllers.WithdrawBalance)
	customerGroup.POST("/transfer", controllers.Transfer)
	customerGroup.GET("/history", controllers.HistoryCustomer)

	restAPI.Run(":8080")

}
