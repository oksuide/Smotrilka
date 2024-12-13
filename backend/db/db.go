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

// Подключаемся к базе данных
func Connect(connStr string) {
	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	// Проверяем соединение с базой данных
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

// InitTables инициализирует таблицы
func InitTables() {
	// Создание таблицы пользователей
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		   username VARCHAR(50) NOT NULL,
		password VARCHAR(50) NOT NULL,
		room_id UUID
	);
    CREATE TABLE IF NOT EXISTS rooms (
    	id UUID PRIMARY KEY,
    	name VARCHAR(50) NOT NULL,
		password VARCHAR(50) NOT NULL,
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		creator INT REFERENCES users(id)
	);
	CREATE TABLE IF NOT EXISTS room_events (
	    id SERIAL PRIMARY KEY,
	    room_id UUID REFERENCES rooms(id),
	    user_id INT REFERENCES users(id),
	    event_type TEXT,
	    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	err := DB.Exec(query).Error
	if err != nil {
		log.Fatal("Error initializing tables:", err)
	}
}
