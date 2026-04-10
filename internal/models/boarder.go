package models

import (
	"errors"
	"time"
)

const (
	StatusBoarderActive     = "active"
	StatusBoarderCheckedOut = "checked_out"
	StatusBoarderPending    = "pending"
)

type Boarder struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FirstName string    `json:"firstName" example:"John"`
	LastName  string    `json:"lastName" example:"Doe"`
	Phone     string    `json:"phone" example:"+1234567890"`
	CreatedAt time.Time `json:"createdAt" example:"2024-01-01T12:00:00Z"`
	UpdatedAt time.Time `json:"updatedAt" example:"2024-01-01T12:00:00Z"`
	RoomID    string    `json:"roomId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status    string    `json:"status" example:"active"`
}

func (b *Boarder) Validate() error {
	if b.FirstName == "" && b.LastName == "" {
		return errors.New("Name is required")
	}
	if b.Phone == "" {
		return errors.New("phone number is required")
	}
	if b.RoomID == "" {
		return errors.New("room ID is required")
	}
	switch b.Status {
	case StatusBoarderActive, StatusBoarderCheckedOut, StatusBoarderPending:
		return nil
	case "":
		return errors.New("status is required")
	default:
		return errors.New("invalid status: must be active, checked_out, or pending")
	}
}
