package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"project-backend/db"
	"project-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ConnectToRoom(c *gin.Context) {
	// Получение ID комнаты из параметров URL
	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	// Проверяем наличие комнаты с таким ID
	var existingRoom models.Room
	if err := db.DB.Where("id = ?", roomID).First(&existingRoom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Комната с таким id не найден
			c.JSON(404, gin.H{"error": "Room not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	// Получение пользователя из токена
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Проверка нахождения пользователя в другой комнате
	if user.RoomID.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already in a room"})
		return
	}

	// Присвоение пользователю `room_id`
	user.RoomID = sql.NullString{String: existingRoom.ID, Valid: true}
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to room"})
		return
	}

	// Увеличение count пользователей в комнате
	existingRoom.UserCount += 1
	if err := db.DB.Save(&existingRoom).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room"})
		return
	}

	if err := LogRoomEvent(existingRoom.ID, user.ID, "UserConnected"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log room event"})
		return
	}
	// Возвращение успешного ответа
	c.JSON(http.StatusOK, gin.H{
		"message": "Connected to room successfully",
		"room":    existingRoom.ID,
		"user":    user,
	})
}

func DisconnectFromRoom(c *gin.Context) {
	// Получение ID комнаты из параметров URL
	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	// Проверяем наличие комнаты с таким ID
	var room models.Room
	if err := db.DB.Where("id = ?", roomID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Комната с таким id не найден
			c.JSON(404, gin.H{"error": "Room not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	// Получение пользователя из токена
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Приведение userID к подходящему типу
	userID, ok := userIDRaw.(int) // Если у тебя UUID, используй string
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Проверка, находится ли пользователь в комнате
	if user.RoomID.String != roomID || !user.RoomID.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not in the specified room"})
		return
	}

	// Начало транзакции для атомарности операций
	tx := db.DB.Begin()

	// Присвоение пользователю пустого `room_id`
	user.RoomID = sql.NullString{String: "", Valid: false}
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disconnect from the room"})
		return
	}

	// Уменьшение `UserCount` в комнате
	if room.UserCount > 0 {
		room.UserCount -= 1
	}
	if err := tx.Save(&room).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room"})
		return
	}

	// Коммит транзакции
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	if err := LogRoomEvent(room.ID, user.ID, "UserDisconnected"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log room event"})
		return
	}
	// Возвращение успешного ответа
	c.JSON(http.StatusOK, gin.H{
		"message": "Disconnected from the room successfully",
		"room":    room.ID,
		"user":    user.ID,
	})
}

// crd Room Endpoints
func CreateRoom(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(400, gin.H{"error": "User ID not found in context"})
		return
	}

	// Проверяем наличие пользователя с таким ID
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
	room.ID = roomInd()
	room.Creator = user.ID
	// Хешируем пароль перед сохранением
	hashedPassword, err := hashPassword(room.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error hashing password"})
		return
	}
	room.Password = hashedPassword
	// Сохраняем комнату в базе данных
	if err := db.DB.Create(&room).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error creating room"})
		log.Printf("error:%v", err)
		return
	}

	if err := LogRoomEvent(room.ID, user.ID, "RoomCreated"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log room event"})
		log.Printf("error:%v", err)
		return
	}
	// Возвращаем успешный ответ с данными о созданной комнате (без пароля)
	c.JSON(201, gin.H{
		"id":      room.ID,
		"name":    room.Name,
		"date":    room.CreatedAt,
		"creator": room.Creator,
	})
}

func GetRoom(c *gin.Context) {
	// Получение ID комнаты
	var roomID models.Connect
	if err := c.ShouldBindJSON(&roomID); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Проверяем наличие комнаты с таким ID
	var existingRoom models.Room
	if err := db.DB.Where("id = ?", roomID.ID).First(&existingRoom).Error; err != nil {
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
		"id":         existingRoom.ID,
		"name":       existingRoom.Name,
		"user_count": existingRoom.UserCount,
		"Creator":    existingRoom.Creator,
	})
}

func DeleteRoom(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	// Проверяем наличие пользователя с таким ID
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

	// Получаем ID комнаты из контекста
	// Получение ID комнаты
	var roomID models.Connect
	if err := c.ShouldBindJSON(&roomID); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Проверяем наличие комнаты с таким ID
	var existingRoom models.Room
	if err := db.DB.Where("id = ?", roomID.ID).First(&existingRoom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Комната с таким id не найден
			c.JSON(404, gin.H{"error": "Room not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	if err := LogRoomEvent(existingRoom.ID, user.ID, "RoomDeleted"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log room event"})
		return
	}

	// Проверяем является пользователь создателем комнаты и удаляем комнату
	if user.ID == existingRoom.Creator {
		if err := db.DB.Delete(&existingRoom).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete room"})
			return
		}
	} else {
		c.JSON(401, gin.H{"error": "The user is not the creator of the room"})
		return
	}

	c.JSON(200, gin.H{"message": "Room deleted successfully"})
}
func roomInd() string {
	return uuid.New().String()
}
