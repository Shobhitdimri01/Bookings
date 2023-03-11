package repository

import (
	"time"

	"github.com/Shobhitdimri01/Bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	GetUserById(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string,int, error)
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(u models.Reservation) error
	DeleteReservation(id int)error
	UpdateProcessedForReservation(id,processed int) error
	CountReservation()(int,int,int)
	CountMonths()([]int,[]int)
	InsertUserData(models.User) error
	EmailCheck(email string)bool
	GetAllAdmins()([]models.User,error)
	GetAdminByID(id int)(models.User,error)
	DeleteAdminByID(id int)error
	UpdateAdminData(u models.User)error
	AllRooms()([]models.Room,error)
	GetRestrictionsForRoomByDate(roomID int,startDate,endDate time.Time)([]models.RoomRestriction,error)
	InsertBlockForRoom(id int, startDate time.Time) error
	DeleteBlockByID(id int) error
}
