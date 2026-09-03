package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(expectedToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ignora verificação para requisições Preflight do CORS
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")

		// Remove o prefixo "Bearer " se existir e faz o trim
		rawToken := strings.TrimPrefix(token, "Bearer ")
		rawToken = strings.TrimSpace(rawToken)

		if rawToken == "" || rawToken != strings.TrimSpace(expectedToken) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Acesso não autorizado: Token inválido ou ausente",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
