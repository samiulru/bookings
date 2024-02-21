package dbrepo

import (
	"errors"
	"github.com/samiulru/bookings/internal/models"
	"time"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// Reservations inserts reservation info to the database
func (m *testDBRepo) InsertReservations(res models.Reservation) (int, error) {
	if res.RoomID == 100 {
		return 0, errors.New("Unable to reserve this room")
	}
	return 1, nil
}

// InsertRoomRestriction inserts room restriction info to the database
func (m *testDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	if res.RoomID == 200 {
		return errors.New("Unable to insert room restriction this room")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID searchs and return if any room is available for a specific date range
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	//Invalid Date pair
	//layout := "02-01-2006"
	//sd := start.Format(layout)
	//ed := end.Format(layout)
	//Invalid Date pair
	//if sd == "01-01-2050" && ed == "02-01-2050" {
	//	return false, nil
	//}
	if roomID == 2 {
		return false, errors.New("No room available")
	}
	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of all availabile rooms for a specific date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room

	layout := "02-01-2006"
	sd := start.Format(layout)
	ed := end.Format(layout)
	//Invalid Date pair
	if sd == "02-01-2050" && ed == "01-01-2050" {
		return rooms, errors.New("Invalid Date pair, end date can't be before the start date")
	}

	//Room Unavailable for this time span
	if sd == "01-01-2050" && ed == "02-01-2050" {
		return rooms, nil
	}
	rooms = append(rooms, models.Room{})
	return rooms, nil
}

// GetRoomByID searchs room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("Non existent room")
	}
	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	var id int
	var hashedPassword string
	return id, hashedPassword, nil
}

// ViewALlReservations returns a slice of all reservations
func (m *testDBRepo) ViewALlReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	var err error
	return reservations, err
}

// ViewNewReservations returns a slice of all reservations
func (m *testDBRepo) ViewNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

// GetReservationByID searches reservation info by ID
func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var reservation models.Reservation
	return reservation, nil
}

// UpdatReservation updates reservation information in the database
func (m *testDBRepo) UpdateReservation(r models.Reservation) error {
	return nil
}

// DeleteReservation deletes reservation information from the database
func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

// UpdateProcessedForReservation updates Processed information in the database
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}

// AllRooms returns a slice of all room
func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	return nil, nil
}

// GetRestrictionsForRoomByID returns a slice of all restrictions for a room within the given date range
func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start_date, end_date time.Time) ([]models.RoomRestriction, error) {

	return nil, nil
}

