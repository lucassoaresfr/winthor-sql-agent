package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(expectedToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" || (token != expectedToken && token != "Bearer "+expectedToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Acesso não autorizado: Token inválido ou ausente"})
			c.Abort()
			return
		}
		c.Next()
	}

}
