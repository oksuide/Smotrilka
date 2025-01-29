package queries

import (
	"project-backend/db"
	"project-backend/models"
)

// Проверяем наличие комнаты с таким ID
func RoomSearch(roomID string) (models.Room, error) {
	var existingRoom models.Room
	if err := db.DB.Where("id = ?", roomID).First(&existingRoom).Error; err != nil {
		return existingRoom, err
	}
	return existingRoom, nil
}
