package main

import (
	"log"
	"net/http"
	"os"
	"project-backend/db"
	"project-backend/handlers"
	"project-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Println("Message received:", string(message))
	}
}

func main() {
	// Подключаемся к .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Подключение к базе данных
	connStr := os.Getenv("ConnStr")
	db.Connect(connStr)
	db.InitTables()

	// Создание роутера
	router := gin.Default()

	// Применяем AuthMiddleware ко всем маршрутам, требующим авторизации
	authorized := router.Group("/api")
	authorized.Use(middleware.AuthMiddleware())
	{
		// Роуты для работы с пользователями
		authorized.GET("/user", handlers.GetUser)
		authorized.PUT("/user", handlers.UpdateUser)
		authorized.DELETE("/users/:id", handlers.DeleteUser)
		authorized.POST("/users", handlers.CreateUser)

		// Роуты для комнат
		rooms := authorized.Group("/rooms")
		{
			rooms.Use(middleware.RoomIDMiddleware()) // Применяем middleware для roomID
			rooms.GET("/:postID", handlers.GetRoom)
			rooms.DELETE("/:postID", handlers.DeleteRoom)
			rooms.POST("/", handlers.CreateRoom)
		}

		// Запуск сервера на порту 8080
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("Server start error: %v", err)
		}
	}
}
