package middleware

import (
	"agendamento-salas/internal/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware intercepta e valida requisições com JWT Bearer Token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretStr := os.Getenv("JWT_SECRET")
		if secretStr == "" {
			secretStr = "supersecretkey_change_in_prod"
		}
		secret := []byte(secretStr)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorizacao ausente"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token invalido, esperado: Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		userID, err := utils.ValidateToken(tokenString, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalido ou expirado", "detalhe": err.Error()})
			c.Abort()
			return
		}

		// Injeta user_id no contexto
		c.Set("user_id", userID)
		c.Next()
	}
}
