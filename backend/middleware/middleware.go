package middleware

import (
	"log"
	"net/http"
	"project-backend/auth"
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
			log.Println("Authorization header is missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Токен должен начинаться с "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			log.Println("Bearer token is missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is missing"})
			c.Abort()
			return
		}

		// log.Println("Authorization header:", authHeader)
		// log.Println("Token string after trimming 'Bearer ':", tokenString)

		// Парсим и проверяем токен
		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return auth.JwtKey, nil
		})

		// Логирование ошибок при парсинге
		if err != nil {
			log.Println("Error parsing token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Логирование валидности токена
		if !token.Valid {
			log.Println("Token is not valid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Сохраняем userID в контексте для дальнейшего использования
		c.Set("userID", claims.UserID)

		// Логирование сохраненного userID
		userID, exists := c.Get("userID")
		if !exists {
			log.Println("User ID not found in context")
		} else {
			log.Println("User ID from context:", userID)
		}

		// Переходим к следующему обработчику
		c.Next()
	}
}

// func RoomIDMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		roomID := c.Param("roomID")
// 		if roomID == "" {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "roomID is required"})
// 			c.Abort()
// 			return
// 		}

// 		// Конвертируем roomID в int
// 		id, err := strconv.Atoi(roomID)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roomID"})
// 			c.Abort()
// 			return
// 		}

// 		// Сохраняем roomID в контексте
// 		c.Set("roomID", id)
// 		c.Next()
// 	}
// }
