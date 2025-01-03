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

// Настройка Upgrader для работы с WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Обработка соединения WebSocket
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Println("Message received:", string(message))
	}
}

func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Подключение к базе данных
	connStr := os.Getenv("ConnStr")
	db.Connect(connStr)
	db.InitTables()

	// Создание роутера и настройка маршрутов
	router := setupRouter()
	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server start error: %v", err)
	}
}

// Настройка маршруты и middleware
func setupRouter() *gin.Engine {
	router := gin.Default()

	// Роут для WebSocket
	router.GET("/ws", gin.WrapH(http.HandlerFunc(handleWebSocket)))

	// Роуты API
	api := router.Group("/api")

	// Роуты без авторизации
	api.GET("/rooms/:id/events", handlers.GetRoomEvents)
	api.POST("/users", handlers.CreateUser)
	api.POST("/login", handlers.Login)
	// Роуты с авторизацией
	authorized := api.Group("")
	authorized.Use(middleware.AuthMiddleware())
	{
		setupUserRoutes(authorized)
		setupRoomRoutes(authorized)
	}

	return router
}

// setupUserRoutes настраивает маршруты для работы с пользователями
func setupUserRoutes(group *gin.RouterGroup) {
	users := group.Group("/users")
	{
		users.GET("/:id", handlers.GetUser)
		users.PUT("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)
	}

}

// setupRoomRoutes настраивает маршруты для работы с комнатами
func setupRoomRoutes(group *gin.RouterGroup) {
	rooms := group.Group("/rooms")
	rooms.Use(middleware.RoomIDMiddleware())
	{
		rooms.POST("/", handlers.CreateRoom)
		rooms.GET("/:id", handlers.GetRoom)
		rooms.DELETE("/:id", handlers.DeleteRoom)
	}
}
