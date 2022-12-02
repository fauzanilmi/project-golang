package controllers

import (
	"fmt"
	"net/http"
	"project/models"
	"project/token"

	"github.com/gin-gonic/gin"
)

func CheckBalance(ctx *gin.Context) {
	auth, shouldReturn := CheckAuth(ctx)
	if shouldReturn {
		return
	}

	c, err := models.GetBalancebyId(auth.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := fmt.Sprintf("Hello, %s! Your balance is $%d", c.Name, c.Balance)

	ctx.JSON(http.StatusOK, gin.H{"message": message})
}

type TopUpInput struct {
	Amount uint `json:"amount" binding:"required"`
}

func TopUpBalance(ctx *gin.Context) {

	auth, shouldReturn := CheckAuth(ctx)
	if shouldReturn {
		return
	}

	var input TopUpInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c, err := models.GetBalancebyId(auth.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var CurrentBalance uint = c.Balance

	tx := models.Transaction{}
	tx.Amount = input.Amount
	tx.Description = "topup"
	tx.CustomerId = auth.UserId
	tx.Balance = CurrentBalance + input.Amount

	_, err = tx.UpdateBalance()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := fmt.Sprintf("Topup success! Your balance is $%d", tx.Balance)

	ctx.JSON(http.StatusOK, gin.H{"message": message})
}

func WithdrawBalance(ctx *gin.Context) {

	auth, shouldReturn := CheckAuth(ctx)
	if shouldReturn {
		return
	}

	var input TopUpInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c, err := models.GetBalancebyId(auth.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var CurrentBalance uint = c.Balance
	if CurrentBalance < input.Amount {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance!"})
		return
	}
	tx := models.Transaction{}
	tx.Amount = input.Amount
	tx.Description = "withdraw"
	tx.CustomerId = auth.UserId
	tx.Balance = CurrentBalance - input.Amount

	_, err = tx.UpdateBalance()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := fmt.Sprintf("Topup success! Your balance is $%d", tx.Balance)

	ctx.JSON(http.StatusOK, gin.H{"message": message})
}

type TransferInput struct {
	Amount   uint   `json:"amount" binding:"required"`
	Username string `json:"username" binding:"required"`
}

func Transfer(ctx *gin.Context) {

	auth, shouldReturn := CheckAuth(ctx)
	if shouldReturn {
		return
	}

	var input TransferInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c, err := models.GetBalancebyId(auth.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c2, err := models.GetBalancebyUsername(input.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if c.Username == c2.Username {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request!"})
		return
	}

	var CurrentBalance uint = c.Balance
	if CurrentBalance < input.Amount {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance!"})
		return
	}

	tx := models.Transaction{}
	tx.Amount = input.Amount
	tx.Description = "transferTo" + c2.Username
	tx.CustomerId = auth.UserId
	tx.Balance = CurrentBalance - input.Amount

	res, err := tx.UpdateBalance()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	CurrentBalance = c2.Balance
	tx2 := models.Transaction{}
	tx2.Amount = input.Amount
	tx2.Description = "transferFrom" + c.Username
	tx2.CustomerId = c2.ID
	tx2.Balance = CurrentBalance + input.Amount

	_, err = tx2.UpdateBalance()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := fmt.Sprintf("Transfer success! Your balance is $%d", res.Balance)

	ctx.JSON(http.StatusOK, gin.H{"message": message})
}

func HistoryCustomer(ctx *gin.Context) {
	auth, shouldReturn := CheckAuth(ctx)
	if shouldReturn {
		return
	}

	tx, err := models.HistoryCustomer(auth.UserId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tx})

}

func CheckAuth(ctx *gin.Context) (*token.AuthDetails, bool) {
	auth, err := token.ExtractTokenAuth(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return nil, true
	}

	checkAuth := models.CheckAuth(auth)
	if checkAuth != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return nil, true
	}
	return auth, false
}
