package handlers

import (
	"log"
	"net/http"
	"project-backend/auth"
	"project-backend/db"
	"project-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var user models.User
	// Привязка JSON-данных к переменной user
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Проверка уникальности username
	var existingUser models.User
	if err := db.DB.Where("username = ?", user.UserName).First(&existingUser).Error; err == nil {
		c.JSON(400, gin.H{"error": "User with this name already exists"})
		return
	}

	// Хеширования пароля
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error hashing password"})
		return
	}
	log.Printf("Hashed password: %s", hashedPassword)
	user.Password = hashedPassword
	// Сохраняем пользователя в базе данных
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(201, gin.H{
		"id":       user.ID,
		"username": user.UserName,
	})
}

func Login(c *gin.Context) {
	var payload models.LoginPayload

	// Проверяем валидность входящих данных
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Ищем пользователя в базе данных
	var user models.User
	if err := db.DB.Where("username = ?", payload.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
		return
	}

	log.Printf("Stored hash: %s", user.Password)
	log.Printf("Plain password: %s", payload.Password)

	// Сравниваем пароли
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		log.Println("Password comparison failed:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Генерируем JWT-токен
	token, err := auth.GenerateJWT(int(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Возвращаем токен клиенту
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
