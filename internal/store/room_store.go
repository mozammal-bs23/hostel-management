package store

import "example.com/hostel-management/internal/models"

type RoomStore interface {
	Create(room models.Room) (models.Room, error)
	GetByID(id string) (models.Room, ErrNotFound)
	List(filters RoomFilters) ([]models.Room, ErrNotFound)
	Update(id string, room models.Room) (models.Room, ErrNotFound)
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

type ErrNotFound struct {
	Message string
}

func (e *ErrNotFound) Error() string {
	return e.Message
}