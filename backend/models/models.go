package models

type Room struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Password  string `json:"-"`
	CreatedAt int    `json:"date"`
	Creator   int    `json:"creator"`
}

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	UserName string `json:"username"`
	Password string `json:"-"`
	RoomID   int    `json:"room_id"`
}

type RoomEvents struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	RoomID    int    `json:"room_id"`
	UserID    int    `json:"user_id"`
	EventType string `json:"event_type"`
	TimeStamp int    `json:"timestamp"`
}
