package main

import (
	"log"

	"dvfucar/internal/config"
	"dvfucar/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Подключаемся к БД
	config.ConnectDatabase()

	r := gin.Default()

	// Настройка CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Роуты API
	api := r.Group("/api")
	{
		api.POST("/users", handlers.GetOrCreateUser)
		api.GET("/users/:tgId/trips", handlers.GetUserTrips)

		api.GET("/trips", handlers.GetAllTrips)
		api.GET("/search_trips", handlers.SearchTrips)
		api.GET("/trips/:id", handlers.GetTripById)
		api.POST("/trips", handlers.CreateTrip)
		api.POST("/trips/:id/book", handlers.BookTrip)

		api.PATCH("/bookings/:id", handlers.UpdateBookingStatus)
	}

	// Роут WebSocket
	r.GET("/ws/trip/:id", handlers.HandleWebSocket)

	log.Println("Server is running on port 8000...")
	r.Run(":8000")
}
