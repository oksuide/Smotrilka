package handlers

import (
	"net/http"
	"project-backend/db"
	"project-backend/models"

	"github.com/gin-gonic/gin"
)

func GetRoomEvents(c *gin.Context) {
	roomID := c.Param("id")

	var events []models.RoomEvent
	if err := db.DB.Where("room_id = ?", roomID).Order("timestamp desc").Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}
func LogRoomEvent(roomID string, userID uint, eventType string) error {
	event := models.RoomEvent{
		RoomID:    roomID,
		UserID:    userID,
		EventType: eventType,
	}
	if err := db.DB.Create(&event).Error; err != nil {
		return err
	}
	return nil
}
