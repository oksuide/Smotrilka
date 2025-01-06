package db

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

// Подключение к базе данных
func Connect(connStr string) {
	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	// Проверка соединение с базой данных
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("Error getting DB instance:", err)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		fmt.Println("Database connection failed:", err)
	} else {
		fmt.Println("Database connection successful!")
	}

}

func InitTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL,
			password TEXT NOT NULL,
			room_id UUID DEFAULT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS rooms (
			id UUID PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			password VARCHAR(255) NOT NULL,
			user_count INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			creator INT REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS room_events (
			id SERIAL PRIMARY KEY,
			room_id UUID NOT NULL,
			user_id INT NOT NULL,
			event_type VARCHAR(50) NOT NULL,
			time_stamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if err := DB.Exec(query).Error; err != nil {
			log.Fatalf("Error creating table: %v\nRequest: %s", err, query)
		}
	}
}
