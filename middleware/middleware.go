package middleware

import (
	"net/http"
	"project/token"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := token.TokenValid(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
