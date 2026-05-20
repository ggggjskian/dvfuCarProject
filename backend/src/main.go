package main

import (
	"log"
	"net/http"
	"strconv"

	"dvfucar/internal/config"
	"dvfucar/internal/handlers"
	"dvfucar/internal/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		api.POST("/users", handlers.GetOrCreateUser)
		api.GET("/users/:tgId/trips", handlers.GetUserTrips)

		api.GET("/trips", handlers.GetAllTrips)
		api.GET("/search_trips", handlers.SearchTrips)
		api.GET("/trips/:id", handlers.GetTripById)

		authorized := api.Group("/")
		authorized.Use(handlers.AuthMiddleware())
		{
			authorized.POST("/trips", createTripWithID)
			authorized.POST("/trips/:id/book", bookTripWithID)
			authorized.PATCH("/bookings/:id", handlers.UpdateBookingStatus)
		}
	}

	r.GET("/ws/trip/:id", handlers.HandleWebSocket)

	log.Println("Server is running on port 8000...")
	r.Run(":8000")
}

func createTripWithID(c *gin.Context) {
	var trip models.Trip
	if err := c.ShouldBindJSON(&trip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Пользователь не авторизован бэкендом"})
		return
	}
	trip.DriverID = userID.(int64)

	if err := config.DB.Create(&trip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, trip)
}

func bookTripWithID(c *gin.Context) {
	tripIDStr := c.Param("id")
	tripID, err := strconv.ParseUint(tripIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID поездки"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Пользователь не авторизован бэкендом"})
		return
	}

	booking := models.Booking{
		TripID:      uint(tripID),
		PassengerID: userID.(int64),
		Status:      "pending",
	}

	if err := config.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, booking)
}
