package repository

import "github.com/Shobhitdimri01/Bookings/internal/models"

type DataBaseRepo interface {
	AllUsers() 						bool
	
	InsertReservation(res models.Reservation)	error
}
