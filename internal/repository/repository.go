package repository

import (
	"github.com/samiulru/bookings/internal/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservations(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
}
