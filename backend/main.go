package main

import (
	"log"
	"net/http"
	"os"
	"project-backend/db"
	"project-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	// Подключение к базе данных
	connStr := os.Getenv("ConnStr")
	db.Connect(connStr)
	db.InitTables()

	// Настройка роутов
	router := routes.SetupRouter()

	// Настройка CORS
	handler := setupCORS(router)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server start error: %v", err)
	}
}

// Функция для настройки CORS
func setupCORS(router *gin.Engine) http.Handler {
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(router)

	return handler
}
