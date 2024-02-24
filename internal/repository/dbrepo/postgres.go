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

// InsertReservations inserts reservation info to the database
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

// AllRooms returns a slice of all room
func (m *postgresDBRepo) AllRooms() ([]models.Room, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []models.Room
	query := `select id, room_name, created_at, updated_at from rooms order by room_name`

	rows, err := m.DB.QueryContext(cntx, query)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var room models.Room
		err = rows.Scan(
			&room.ID,
			&room.RoomName,
			&room.CreatedAt,
			&room.UpdateAt,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
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

// ViewALlReservations returns a slice of all reservations
func (m *postgresDBRepo) ViewALlReservations() ([]models.Reservation, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation
	query := `select r.id, r.first_name, r.last_name, r.email, r.mobile_number, r.start_date, r.end_date,
       r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
		from reservations r 
		left join rooms rm on (r.room_id = rm.id)
		order by r.start_date asc`

	rows, err := m.DB.QueryContext(cntx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation
		err = rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.MobileNumber,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdateAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

// ViewNewReservations returns a slice of all reservations
func (m *postgresDBRepo) ViewNewReservations() ([]models.Reservation, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation
	query := `select r.id, r.first_name, r.last_name, r.email, r.mobile_number, r.start_date, r.end_date,
       r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
		from reservations r 
		left join rooms rm on (r.room_id = rm.id)
		where processed = 0
		order by r.start_date asc`

	rows, err := m.DB.QueryContext(cntx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation
		err = rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.MobileNumber,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdateAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, err
}

// GetReservationByID searches reservation info by ID
func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservation models.Reservation
	query := `select r.id, r.first_name, r.last_name, r.email, r.mobile_number, r.start_date, r.end_date,
       r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
		from reservations r 
		left join rooms rm on (r.room_id = rm.id)
		where r.id = $1
`

	rows := m.DB.QueryRowContext(cntx, query, id)
	err := rows.Scan(
		&reservation.ID,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Email,
		&reservation.MobileNumber,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.RoomID,
		&reservation.CreatedAt,
		&reservation.UpdateAt,
		&reservation.Processed,
		&reservation.Room.ID,
		&reservation.Room.RoomName,
	)

	if err != nil {
		return reservation, err
	}
	return reservation, nil
}

// UpdateReservation updates reservation information in the database
func (m *postgresDBRepo) UpdateReservation(r models.Reservation) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update reservations
			set first_name = $1 , last_name = $2 , email = $3 , mobile_number = $4 , updated_at = $5
			where id = $6`

	_, err := m.DB.ExecContext(cntx, query,
		r.FirstName,
		r.LastName,
		r.Email,
		r.MobileNumber,
		time.Now(),
		r.ID,
	)

	return err
}

// UpdateProcessedForReservation updates Processed information in the database
func (m *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update reservations set processed = $1 where id = $2`

	_, err := m.DB.ExecContext(cntx, query, processed, id)

	return err
}

// DeleteReservation deletes reservation information from the database
func (m *postgresDBRepo) DeleteReservation(id int) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `delete from reservations where id = $1`

	_, err := m.DB.ExecContext(cntx, query, id)

	return err
}

// GetRestrictionsForRoomByID returns a slice of all restrictions for a room within the given date range
func (m *postgresDBRepo) GetRestrictionsForRoomByDate(roomID int, start_date, end_date time.Time) ([]models.RoomRestriction, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var roomRestrictions []models.RoomRestriction

		query := `select id, start_date, end_date, room_id, coalesce(reservation_id, 0), restriction_id, created_at, updated_at 
		from room_restrictions where start_date <= $1 and end_date > $2 and room_id = $3`

	rows, err := m.DB.QueryContext(cntx, query, end_date, start_date, roomID)
	if err != nil {
		return roomRestrictions, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.RoomRestriction
		err = rows.Scan(
			&r.ID,
			&r.StartDate,
			&r.EndDate,
			&r.RoomID,
			&r.ReservationID,
			&r.RestrictionId,
			&r.CreatedAt,
			&r.UpdateAt,
		)

		if err != nil {
			return roomRestrictions, err
		}
		roomRestrictions = append(roomRestrictions, r)
	}

	if err = rows.Err(); err != nil {
		return roomRestrictions, err
	}
	return roomRestrictions, nil
}

// InsertBlockForRoom inserts room restriction due to room-upgradation or maintainacne 
func (m *postgresDBRepo) InsertBlockForRoom(roomID int, start_date, end_date time.Time) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, restriction_id, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.ExecContext(cntx, stmt,
		start_date,
		end_date,
		roomID,
		2,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteBlockForRoom deletes room restriction due to completed room-upgradation or maintainacne 
func (m *postgresDBRepo) DeleteBlockForRoom(id int) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `delete from room_restrictions where id = $1`

	_, err := m.DB.ExecContext(cntx, stmt, id)

	if err != nil {
		return err
	}

	return nil
}