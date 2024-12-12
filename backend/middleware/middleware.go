package middleware

import (
	"net/http"
	"project-backend/auth"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Middleware для проверки JWT токена
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Токен должен начинаться с "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is missing"})
			c.Abort()
			return
		}

		// Парсим и проверяем токен
		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return auth.JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Сохраняем userID в контексте для дальнейшего использования
		c.Set("userID", claims.UserID)

		// Переходим к следующему обработчику
		c.Next()
	}
}

func RoomIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("roomID")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "roomID is required"})
			c.Abort()
			return
		}

		// Конвертируем roomID в int
		id, err := strconv.Atoi(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roomID"})
			c.Abort()
			return
		}

		// Сохраняем roomID в контексте
		c.Set("roomID", id)
		c.Next()
	}
}
