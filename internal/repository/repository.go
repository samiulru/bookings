package repository

import "github.com/samiulru/bookings/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservations(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
}
