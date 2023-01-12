package models

import (
	"time"
)

// // Reservation holds reservation data
// type Reservation struct {
// 	FirstName  		string
// 	LastName 		 string
// 	Email     			 string
// 	Phone    			 string
// }

//Structs for Database :-

// User is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room is the room model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restrictions is the restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation is reservations model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
	Processed int
}

// RoomRestriction is RoomRestrictions model
type RoomRestriction struct {
	ID              int
	RestrictionName string
	RoomID          int
	Room            Room
	Reservation     Reservation
	Restrictions    Restriction
	ReservationID   int
	RestrictionID   int
	StartDate       time.Time
	EndDate         time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Maildata holds an email message
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}
