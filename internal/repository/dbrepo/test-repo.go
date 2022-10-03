package dbrepo


import (
	"time"
	"github.com/Shobhitdimri01/Bookings/internal/models"
)

func (m *testDbRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDbRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDbRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	return nil
}

//Will check the date and See whether there is availability in room for roomId --> returns True or false --> if no Availability
func (m *testDbRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

//SearchAvailabilityforAll rooms returns a slice of available rooms if any , for given date range
func (m *testDbRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil

}

//GetRoomByID will get room name with its defined Id
func (m *testDbRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	return room, nil

}
