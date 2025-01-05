package models

import (
	"database/sql"
	"time"
)

type Room struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Password  string    `json:"password" gorm:"not null"`
	UserCount int       `json:"user_count" gorm:"default:0"`
	CreatedAt time.Time `json:"date" gorm:"not null"`
	Creator   uint      `json:"creator" gorm:"not null"`
}

type User struct {
	ID       uint           `json:"id" gorm:"primaryKey"`
	UserName string         `json:"username" gorm:"column:username;not null"`
	Password string         `json:"password" gorm:"not null"`
	RoomID   sql.NullString `json:"room_id" gorm:"default:null"`
}

type RoomEvent struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RoomID    string    `json:"room_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	EventType string    `json:"event_type" gorm:"not null" gorm:"size:50"`
	TimeStamp time.Time `json:"timestamp" gorm:"default:current_timestamp"`
}
type LoginPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
