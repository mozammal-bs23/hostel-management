package models

import (
	"errors"
	"time"
)

// Room status constants
const (
	StatusAvailable   = "available"
	StatusOccupied    = "occupied"
	StatusMaintenance = "maintenance"
)

// Room represents a hostel room entity
type Room struct {
	ID          string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string    `json:"name" example:"Room 101"`
	Capacity    int       `json:"capacity" example:"2"`
	Status      string    `json:"status" example:"available"`
	RentalPrice float64   `json:"rentalPrice" example:"100.00"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Validate checks if the room data is correct
func (r *Room) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Capacity <= 0 {
		return errors.New("capacity must be at least 1")
	}
	if r.RentalPrice < 0 {
		return errors.New("rental price cannot be negative")
	}
	switch r.Status {
	case StatusAvailable, StatusOccupied, StatusMaintenance:
		return nil
	case "":
		return errors.New("status is required")
	default:
		return errors.New("invalid status: must be available, occupied, or maintenance")
	}
}
