package queries

import (
	"project-backend/db"
	"project-backend/models"
)

func UserSearch(userID any) (models.User, error) {
	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		return user, err
	}
	return user, nil
}
