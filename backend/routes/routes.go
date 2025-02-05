package routes

import (
	"net/http"
	"project-backend/handlers"
	"project-backend/middleware"
	"project-backend/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Роут для WebSocket
	router.GET("/ws", gin.WrapH(http.HandlerFunc(websocket.HandleWebSocket)))

	// Роуты API
	api := router.Group("/api")

	// Роуты без авторизации
	api.GET("/rooms/events", handlers.GetRoomEvents)
	api.POST("/register", handlers.CreateUser)
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

func setupUserRoutes(group *gin.RouterGroup) {
	users := group.Group("/users")
	{
		users.GET("/me", handlers.GetMe)
		users.GET("/:id", handlers.GetUser)
		users.PUT("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)
	}
}

func setupRoomRoutes(group *gin.RouterGroup) {
	rooms := group.Group("/rooms")
	{
		rooms.POST("/", handlers.CreateRoom)
		rooms.GET("/my", handlers.GetUserRooms)
		rooms.GET("/:id", handlers.GetRoom)
		rooms.DELETE("/:id", handlers.DeleteRoom)
		rooms.POST("/connect/:id", handlers.ConnectToRoom)
		rooms.POST("/disconnect/:id", handlers.DisconnectFromRoom)
	}
}
