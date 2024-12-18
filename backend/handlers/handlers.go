package handlers

import (
	"errors"
	"project-backend/db"
	"project-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(c *gin.Context) {
	var user models.User
	// Привязываем JSON-данные из тела запроса к переменной user
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Проверяем уникальность username, чтобы не создать дублирующегося пользователя
	var existingUser models.User
	if err := db.DB.Where("username = ?", user.UserName).First(&existingUser).Error; err == nil {
		c.JSON(400, gin.H{"error": "User with this name already exists"})
		return
	}

	// Хешируем пароль перед сохранением
	_, err := hashPassword(&user.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error hashing password"})
		return
	}

	// Сохраняем пользователя в базе данных
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error creating user"})
		return
	}

	// Возвращаем успешный ответ с данными о созданном пользователе (без пароля)
	c.JSON(201, gin.H{
		"id":       user.ID,
		"username": user.UserName,
	})
}
func GetUser(c *gin.Context) {
	// Получаем ID пользователя из параметров URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
	// Проверяем наличие пользователя с таким ID
	var existingUser models.User
	if err := db.DB.Where("id = ?", id).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Пользователь с таким id не найден
			c.JSON(404, gin.H{"error": "User not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
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

func CreateRoom(c *gin.Context) {
	var room models.Room
	// Привязываем JSON-данные из тела запроса к переменной room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Проверяем уникальность room.name, чтобы не создать дублирующегося пользователя
	var existingRoom models.Room
	if err := db.DB.Where("name = ?", room.Name).First(&existingRoom).Error; err == nil {
		c.JSON(400, gin.H{"error": "Room with this name already exists"})
		return
	}

	// Хешируем пароль перед сохранением
	_, err := hashPassword(&room.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error hashing password"})
		return
	}

	// Сохраняем комнату в базе данных
	if err := db.DB.Create(&room).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error creating room"})
		return
	}

	// Возвращаем успешный ответ с данными о созданной комнате (без пароля)
	c.JSON(201, gin.H{
		"id":   roomInd(),
		"name": room.Name,
		"date": room.CreatedAt,
	})
}

func GetRoom(c *gin.Context) {
	// Получаем ID комнаты из параметров URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid room ID"})
		return
	}
	// Проверяем наличие комнаты с таким ID
	var existingRoom models.Room
	if err := db.DB.Where("id = ?", id).First(&existingRoom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Комната с таким id не найден
			c.JSON(404, gin.H{"error": "Room not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(200, gin.H{
		"id":   existingRoom.ID,
		"name": existingRoom.Name,
	})
}

func DeleteRoom(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	// Ищем комнату в базе данных по ID
	var room models.Room
	if err := db.DB.Where("id = ?", room.ID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Комната с таким id не найдена
			c.JSON(404, gin.H{"error": "Room not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	// Проверяем является пользователь создателем комнаты и удаляем комнату
	if userID == room.Creator {
		if err := db.DB.Delete(&room).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete room"})
			return
		}
	} else {
		c.JSON(401, gin.H{"error": "The user is not the creator of the room"})
		return
	}

	c.JSON(200, gin.H{"message": "Room deleted successfully"})
}

// Шифрование пароля
func hashPassword(password *string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), 14)
	return string(bytes), err
}

// Генерация уникального идентификатора для комнаты
func roomInd() string {
	return uuid.New().String()
}
