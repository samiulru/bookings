package dbrepo

import (
	"context"
	"errors"
	"github.com/samiulru/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// Reservations inserts reservation info to the database
func (m *postgresDBRepo) InsertReservations(res models.Reservation) (int, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int
	query := `insert into reservations (first_name, last_name, email, mobile_number, start_date, end_date, room_id, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(cntx, query,
		res.FirstName,
		res.LastName,
		res.Email,
		res.MobileNumber,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts room restriction info to the database
func (m *postgresDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(cntx, stmt,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.ReservationID,
		res.RestrictionId,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByRoomID searchs and return if any room is available for a specific date range
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var numRows int
	query := `
			select 
				count(id)
			from
				room_restrictions
			where
				room_id = $1
				and $2 < end_date and $3 > start_date;`

	err := m.DB.QueryRowContext(cntx, query, roomID, start, end).Scan(&numRows)

	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of all availabile rooms for a specific date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []models.Room
	query := `
			select 
				r.id, r.room_name
			from
				rooms r
			where r.id not in
			(select room_id from room_restrictions rr where  $1 < rr.end_date and $2 > rr.start_date);
			`

	rows, err := m.DB.QueryContext(cntx, query, start, end)

	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetUserByID searches user by ID
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var u models.User
	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
		from users where id = $1`

	row := m.DB.QueryRowContext(cntx, query, id)
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdateAt,
	)

	if err != nil {
		return u, err
	}
	return u, nil
}

// UpdateUser updates a user in the database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update users
			set first_name = $1 , last_name = $2 , email = $3 , password = $4 , access_level = $5, updated_at = $6
			where id = $7`

	_, err := m.DB.ExecContext(cntx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Password,
		u.AccessLevel,
		time.Now(),
		u.ID,
	)

	return err
}

// Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(cntx, "select id, password from users where email = $1", email)

	err := row.Scan(&id, &hashedPassword)

	if err != nil {
		return id, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// GetRoomByID searches room by ID
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var room models.Room
	query := `select id, room_name, created_at, updated_at from rooms where id = $1`

	row := m.DB.QueryRowContext(cntx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdateAt,
	)

	if err != nil {
		return room, err
	}
	return room, nil
}
