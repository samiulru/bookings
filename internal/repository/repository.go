package repository

import (
	"time"

	"github.com/samiulru/bookings/internal/models"
)

type DatabaseRepo interface {

	
	InsertReservations(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	InsertBlockForRoom(roomID int, start_date, end_date time.Time) error

	DeleteBlockForRoom(id int) error
	DeleteReservation(id int) error

	InsertRoom(room models.Room) (int, error)
	AllRooms() ([]models.Room, error)
	AllRoomsDetails() ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	DeleteRoomByID(id int) error
	UpdateRoom(room models.Room) error
	GetRestrictionsForRoomByDate(roomID int, start_date, end_date time.Time) ([]models.RoomRestriction, error)


	GetReservationByID(id int) (models.Reservation, error)
	ViewALlReservations() ([]models.Reservation, error)
	ViewNewReservations() ([]models.Reservation, error)
	UpdateReservation(r models.Reservation) error
	UpdateProcessedForReservation(id, processed int) error

	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	
	Authenticate(email, testPassword string) (int, string, error)
	AllUsers() bool
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	

}
