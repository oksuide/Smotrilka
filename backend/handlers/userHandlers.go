package handlers

import (
	"errors"
	"project-backend/db"
	"project-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// crud User Endpoints
func GetUser(c *gin.Context) {
	// Получение ID пользователя
	var UserID models.Connect
	if err := c.ShouldBindJSON(&UserID); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Проверяем наличие пользователя с таким ID
	var existingUser models.User
	if err := db.DB.Where("id = ?", UserID.ID).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Пользователь с таким id не найден
			c.JSON(404, gin.H{"error": "User not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	// Возвращаем данные о пользователе
	c.JSON(200, gin.H{
		"id":       existingUser.ID,
		"username": existingUser.UserName,
	})
}

func UpdateUser(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	// Привязываем JSON с изменениями к структуре
	var updateData struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	// Ищем пользователя в базе данных по ID
	var user models.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Пользователь с таким id не найден
			c.JSON(404, gin.H{"error": "User not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	// Обновляем только те поля, которые присутствуют в запросе
	updates := make(map[string]interface{})
	if updateData.Username != "" {
		updates["username"] = updateData.Username
	}
	if updateData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}
		updates["password"] = string(hashedPassword)
	}

	// Выполняем обновление в базе данных
	if len(updates) > 0 {
		if err := db.DB.Model(&user).Updates(updates).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update user"})
			return
		}
	}

	// Возвращаем обновленную информацию о пользователе
	c.JSON(200, gin.H{
		"id":       user.ID,
		"username": user.UserName,
	})
}

func DeleteUser(c *gin.Context) {
	// Логика для удаления пользователя
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	// Ищем пользователя в базе данных по ID
	var user models.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Пользователь с таким id не найден
			c.JSON(404, gin.H{"error": "User not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	// Удаляем пользователя
	if err := db.DB.Delete(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

func GetMe(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	// Проверяем наличие пользователя с таким ID
	var existingUser models.User
	if err := db.DB.Where("id = ?", userID).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Пользователь с таким id не найден
			c.JSON(404, gin.H{"error": "User not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	// Возвращаем данные о пользователе
	c.JSON(200, gin.H{
		"id":       existingUser.ID,
		"username": existingUser.UserName,
	})
}
