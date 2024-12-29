package models

import "database/sql"

type Room struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Password  string `json:"-"`
	UserCount int    `json:"user_count"`
	CreatedAt int    `json:"date"`
	Creator   int    `json:"creator"`
}

type User struct {
	ID       int            `json:"id" gorm:"primaryKey"`
	UserName string         `json:"username"`
	Password string         `json:"-"`
	RoomID   sql.NullString `json:"room_id"`
}

type RoomEvents struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	RoomID    string `json:"room_id"`
	UserID    int    `json:"user_id"`
	EventType string `json:"event_type"`
	TimeStamp int    `json:"timestamp"`
}
