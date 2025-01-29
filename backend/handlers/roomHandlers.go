package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"project-backend/db"
	"project-backend/db/queries"
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
	existingRoom, err := queries.RoomSearch(roomID)
	if err != nil {
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

	// Проверка существования пользователя
	user, err := queries.UserSearch(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Пользователь с таким id не найден
			c.JSON(404, gin.H{"error": "User not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
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
	room, err := queries.RoomSearch(roomID)
	if err != nil {
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

	// Проверка существования пользователя
	user, err := queries.UserSearch(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Пользователь с таким id не найден
			c.JSON(404, gin.H{"error": "User not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
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

func CreateRoom(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(400, gin.H{"error": "User ID not found in context"})
		return
	}

	// Проверка существования пользователя
	user, err := queries.UserSearch(userID)
	if err != nil {
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

	// Проверяем уникальность room.name, чтобы не создать дублирующюся комнату
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
	go func() {
		if err := LogRoomEvent(room.ID, user.ID, "RoomCreated"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log room event"})
			log.Printf("error:%v", err)
			return
		}
	}()

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
	var room models.Connect
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Проверяем наличие комнаты с таким ID
	existingRoom, err := queries.RoomSearch(room.ID)
	if err != nil {
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

	// Проверка существования пользователя
	user, err := queries.UserSearch(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Пользователь с таким id не найден
			c.JSON(404, gin.H{"error": "User not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	// Получение ID комнаты
	var room models.Connect
	if err := c.ShouldBindJSON(&room.ID); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Проверяем наличие комнаты с таким ID
	existingRoom, err := queries.RoomSearch(room.ID)
	if err != nil {
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
