package models

import (
	"time"
)

// User defines the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdateAt    time.Time
}

// Room defines the room model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdateAt  time.Time
}

// Restriction defines the restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdateAt        time.Time
}

// Reservation defines the reservations model
type Reservation struct {
	ID           int
	FirstName    string
	LastName     string
	Email        string
	MobileNumber string
	StartDate    time.Time
	EndDate      time.Time
	RoomID       int
	Room         Room
	CreatedAt    time.Time
	UpdateAt     time.Time
	Processed    int
}

// RoomRestriction defines the restriction model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionId int
	CreatedAt     time.Time
	UpdateAt      time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}

// MailData holds mail messages info
type MailData struct {
	From     string
	To       string
	Subject  string
	Content  string
	Template string
}
