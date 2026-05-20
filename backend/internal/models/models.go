package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TelegramID   *int64    `gorm:"uniqueIndex" json:"telegram_id"`
	Phone        *string   `gorm:"uniqueIndex" json:"phone,omitempty"`
	PasswordHash *string   `json:"-"`
	Username     string    `json:"username"` // сюда подставим username или имя из формы
	CreatedAt    time.Time `json:"created_at"`
}

type Trip struct {
	ID                     uint      `gorm:"primaryKey" json:"id"`
	CreatedAt              time.Time `json:"created_at"`
	DriverID               int64     `json:"driver_tg_id"`
	TripType               string    `json:"trip_type"`
	Point                  string    `json:"point"`
	PointLat               float64   `json:"point_lat"`
	PointLon               float64   `json:"point_lon"`
	DepartureTime          time.Time `json:"departure_time"`
	SeatsTotal             int       `json:"seats_total"`
	Price                  int       `json:"price"`
	Comment                string    `json:"comment"`
	MaxDeviationKm         int       `json:"max_deviation_km"`
	TimeFlexibilityMinutes int       `json:"time_flexibility_minutes"`
	Bookings               []Booking `json:"bookings"`
}

type Booking struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	TripID      uint   `json:"trip_id"`
	PassengerID int64  `json:"passenger_tg_id"`
	Status      string `json:"status" gorm:"default:'pending'"`
}
