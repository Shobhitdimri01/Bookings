package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/Shobhitdimri01/Bookings/internal/models"
)

func (m *testDbRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDbRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2, then fail; otherwise, pass
	if res.RoomID == 2 {
		return 0, errors.New("some error)")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDbRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID, and false if no availability
func (m *testDbRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	// set up a test time
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// this is our test to fail the query -- specify 2060-01-01 as start
	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return false, errors.New("some error")
	}

	// if the start date is after 2049-12-31, then return false,
	// indicating no availability;
	if start.After(t) {
		return false, nil
	}

	// otherwise, we have availability
	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *testDbRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room

	// if the start date is after 2049-12-31, then return empty slice,
	// indicating no rooms are available;
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return rooms, errors.New("some error")
	}

	if start.After(t) {
		return rooms, nil
	}

	// otherwise, put an entry into the slice, indicating that some room is
	// available for search dates
	room := models.Room{
		ID: 1,
	}
	rooms = append(rooms, room)

	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *testDbRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("some error")
	}
	return room, nil
}

func (m *testDbRepo) GetUserById(id int) (models.User, error) {
	var u models.User

	return u, nil
}

func (m *testDbRepo) UpdateUser(u models.User) error {

	return nil
}

func (m *testDbRepo) Authenticate(email, testPassword string) (int, string,int, error) {
	return 1, "",1, nil
}
func (m *testDbRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

func (m *testDbRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}
func (m *testDbRepo) GetReservationByID(id int) (models.Reservation, error) {
	var res models.Reservation
	return res, nil
}

func (m *testDbRepo)UpdateReservation(u models.Reservation) error{
	return nil
}
func (m *testDbRepo)DeleteReservation(id int)error{
	return nil
}
func (m *testDbRepo)UpdateProcessedForReservation(id,processed int) error{
	return nil
}
func (m *testDbRepo)CountReservation()(int,int,int){
	return 1,1,1
}
func (m *testDbRepo)CountMonths()([]int,[]int){
	return nil,nil
}
func (m *testDbRepo)InsertUserData(r models.User) error{
	return nil  
}
func(m *testDbRepo)EmailCheck(email string)bool{
	return true
}
func (m *testDbRepo)GetAllAdmins()([]models.User,error){
	return nil,nil
}
func (m *testDbRepo)GetAdminByID(id int)(models.User,error){
	var admin models.User
	return admin,nil
}
func (m *testDbRepo)DeleteAdminByID(id int)error{
	return nil
}
func(m *testDbRepo)UpdateAdminData(u models.User)error{
	return nil
}
func(m *testDbRepo)AllRooms()([]models.Room,error){
	return nil,nil
}
func(m *testDbRepo)GetRestrictionsForRoomByDate(roomID int,startDate,endDate time.Time)([]models.RoomRestriction,error){
	return nil,nil
}
func(m *testDbRepo)InsertBlockForRoom(id int, startDate time.Time) error{
	return nil
}
func(m *testDbRepo)DeleteBlockByID(id int) error{
	return nil
}