package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hugaojanuario/sentinel/internal/services"
	"github.com/hugaojanuario/sentinel/pkg/config"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token ausente"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenStr, &services.Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.LoadDotEnv().JWT_SECRET), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*services.Claims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			c.Abort()
			return
		}

		c.Set("userId", claims.UserId)
		c.Next()
	}
}
