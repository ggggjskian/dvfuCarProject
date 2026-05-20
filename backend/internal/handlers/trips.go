package handlers

import (
	"dvfucar/internal/config"
	"dvfucar/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateTrip(c *gin.Context) {
	var trip models.Trip
	if err := c.ShouldBindJSON(&trip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}

	tgID, _ := strconv.ParseInt(c.Query("driver_tg_id"), 10, 64)
	trip.DriverID = tgID

	if err := config.DB.Create(&trip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trip"})
		return
	}
	c.JSON(http.StatusCreated, trip)
}

func GetAllTrips(c *gin.Context) {
	var trips []models.Trip
	config.DB.Preload("Bookings").Order("departure_time asc").Find(&trips)
	c.JSON(http.StatusOK, trips)
}

func SearchTrips(c *gin.Context) {
	var trips []models.Trip
	query := config.DB.Model(&models.Trip{})

	if tripType := c.Query("trip_type"); tripType != "" {
		query = query.Where("trip_type = ?", tripType)
	}

	query.Preload("Bookings").Order("departure_time asc").Find(&trips)
	c.JSON(http.StatusOK, trips)
}

func GetTripById(c *gin.Context) {
	id := c.Param("id")
	var trip models.Trip
	if err := config.DB.Preload("Bookings").First(&trip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}
	c.JSON(http.StatusOK, trip)
}

func BookTrip(c *gin.Context) {
	tripID, _ := strconv.Atoi(c.Param("id"))
	pID, _ := strconv.ParseInt(c.Query("passenger_tg_id"), 10, 64)

	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		// Игнорируем ошибку, так как body может быть пустым
	}

	booking.TripID = uint(tripID)
	booking.PassengerID = pID
	booking.Status = "pending"

	if err := config.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to book trip"})
		return
	}
	c.JSON(http.StatusCreated, booking)
}
