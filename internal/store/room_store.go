package store

import (
	"errors"

	"example.com/hostel-management/internal/models"
)

type RoomStore interface {
	Create(room models.Room) (models.Room, error)
	GetByID(id string) (models.Room, error)
	List(filters RoomFilters) ([]models.Room, error)
	Update(id string, room models.Room) (models.Room, error)
	Delete(id string) error
}

type RoomFilters struct {
	Status *string
	Limit int
	Offset int
}

func (r *RoomFilters) IsValid() bool {
	return r.Limit > 0 && r.Offset >= 0
}

var ErrNotFound = errors.New("not found")
