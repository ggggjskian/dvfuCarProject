package handlers

import (
	"dvfucar/internal/config"
	"dvfucar/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOrCreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Where(models.User{TelegramID: user.TelegramID}).FirstOrCreate(&user)
	c.JSON(http.StatusOK, user)
}

func GetUserTrips(c *gin.Context) {
	tgID, _ := strconv.ParseInt(c.Param("tgId"), 10, 64)
	var trips []models.Trip

	config.DB.Preload("Bookings").
		Where("driver_tg_id = ?", tgID).
		Or("id IN (SELECT trip_id FROM bookings WHERE passenger_tg_id = ?)", tgID).
		Find(&trips)

	c.JSON(http.StatusOK, trips)
}
