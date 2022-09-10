package repository

import "github.com/Shobhitdimri01/Bookings/internal/models"

type DatabaseRepo interface {
	AllUsers() 						bool
	
	InsertReservation(res models.Reservation)	(int,error)
	InsertRoomRestriction(r models.RoomRestriction) error
}
