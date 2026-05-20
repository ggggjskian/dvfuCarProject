package handlers

import (
	"dvfucar/internal/config"
	"dvfucar/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateBookingStatus(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&models.Booking{}).Where("id = ?", id).Update("status", input.Status)
	c.JSON(http.StatusOK, gin.H{"status": "Status updated successfully"})
}
