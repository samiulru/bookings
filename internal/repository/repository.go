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

	AllRooms() ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	GetUserByID(id int) (models.User, error)

	Authenticate(email, testPassword string) (int, string, error)
	UpdateUser(u models.User) error
	UpdateReservation(r models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedForReservation(id, processed int) error

	ViewALlReservations() ([]models.Reservation, error)
	ViewNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)

	GetRestrictionsForRoomByDate(roomID int, start_date, end_date time.Time) ([]models.RoomRestriction, error)
}
