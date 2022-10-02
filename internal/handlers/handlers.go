package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	//"log"
	"strconv"
	"time"

	//"log"
	"net/http"

	"github.com/Shobhitdimri01/Bookings/internal/config"
	"github.com/Shobhitdimri01/Bookings/internal/driver"
	"github.com/Shobhitdimri01/Bookings/internal/forms"
	"github.com/Shobhitdimri01/Bookings/internal/helpers"
	"github.com/Shobhitdimri01/Bookings/internal/models"
	"github.com/Shobhitdimri01/Bookings/internal/render"
	"github.com/Shobhitdimri01/Bookings/internal/repository"
	"github.com/Shobhitdimri01/Bookings/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
)

//Repo is repository used by handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

//New Handlers sets the repository for handler
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
//(m *Repository) is the reciever function due to which all of the function is linked with AppConfig
// Home is the handler for the home page
// Home is the handler for the home page
// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.html", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	// stringMap := make(map[string]string)
	// stringMap["test"] = "Hello, again"

	// remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	// stringMap["remote_ip"] = remoteIP

	// send data to the template
	render.Template(w, r, "about.html", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//Extracting roomName from Database
	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, errors.New("can't  get reservation from session"))
		return
	}

	res.Room.RoomName = room.RoomName
	m.App.InfoLog.Println("ROOM:", res.Room.RoomName)
	m.App.Session.Put(r.Context(), "reservation", res)
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

//Post Reservation handles the posting and validation of form
// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("can't get session"))
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//*******************This needs to be remain commented***************************
	// sd := r.Form.Get("start_date")
	// ed := r.Form.Get("end_date")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700

	// layout := "2006-01-02"

	// startDate, err := time.Parse(layout, sd)
	// if err != nil {
	// 	m.App.Session.Put(r.Context(), "error", "can't parse start date")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	// endDate, err := time.Parse(layout, ed)
	// if err != nil {
	// 	m.App.Session.Put(r.Context(), "error", "can't get parse end date")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	// roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	// if err != nil {
	// 	m.App.Session.Put(r.Context(), "error", "invalid data!")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	//Updating my Summary-page
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	// reservation := models.Reservation{
	// 	FirstName: r.Form.Get("first_name"),
	// 	LastName:  r.Form.Get("last_name"),
	// 	Phone:     r.Form.Get("phone"),
	// 	Email:     r.Form.Get("email"),
	// 	// StartDate: startDate,
	// 	// EndDate:   endDate,
	// 	RoomID: roomID,
	// }

	form := forms.New(r.PostForm)

	//minimum requirement specified
	form.Required("first_name", "last_name", "email")
	form.Minlength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		// http.Error(w, "my own error message", http.StatusSeeOther)
		render.Template(w, r, "make-reservation.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		//helpers.ServerError(w, err)
		fmt.Println("My error -->", err.Error())
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// Generals renders the room page
func (m *Repository) Deluxe(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "deluxe.html", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Suite(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "suite.html", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.html", &models.TemplateData{})
}

// PostAvailability renders the search availability page
// PostAvailability handles post
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {

		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	//For terminal purpose//////////////////////////////////////////////////////
	for _, i := range rooms {
		m.App.InfoLog.Println("ROOM:", i.ID, " - ", i.RoomName)

	}
	m.App.InfoLog.Println("---------------------------------------------------")
	////////////////////////////////////////////////////////////////////////////////////
	if len(rooms) == 0 {
		m.App.InfoLog.Println("Sorry !!!     Rooms fully Booked")
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)
	render.Template(w, r, "choose_room.html", &models.TemplateData{
		Data: data,
	})

	//w.Write([]byte(fmt.Sprintf("start date is %s and end is %s", start, end)))
}

type jsonresponse struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

//AvailabilityJson handles request for availability and sends JSON as response
func (m *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startdate, err := time.Parse(layout, sd)
	if err != nil {
		// log.Println(err.Error())
		helpers.ServerError(w, err)
	}
	enddate, err := time.Parse(layout, ed)
	if err != nil {
		// log.Println(err.Error())
		helpers.ServerError(w, err)
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		// log.Println(err.Error())
		helpers.ServerError(w, err)
	}
	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startdate, enddate, roomID)
	if err != nil {
		// log.Println(err.Error())
		helpers.ServerError(w, err)
	}

	resp := jsonresponse{
		Ok:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		// log.Println(err.Error())
		helpers.ServerError(w, err)
	}

	// log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.html", &models.TemplateData{})
}

//Reservation summary displays the user reservation details 
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	//Pulling the data out from Session
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	render.Template(w, r, "reservation-summary.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}


//Choose room displays lists tha available room to the user
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {

	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Get(r.Context(), "reservation")
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}
	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}


//Bookroom takes url parameter build sessional variable and takes user to make reservation screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	//id , s , e
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd:= r.URL.Query().Get("s")
	ed:= r.URL.Query().Get("e")
	log.Println(roomID)

	layout := "2006-01-02"
	startdate, err := time.Parse(layout, sd)
	if err != nil {
		// log.Println(err.Error())
		helpers.ServerError(w, err)
	}
	enddate, err := time.Parse(layout, ed)
	if err != nil {
		// log.Println(err.Error())
		helpers.ServerError(w, err)
	}

	var res models.Reservation
	res.RoomID = roomID
	res.StartDate = startdate
	res.EndDate = enddate

	room , _ :=m.DB.GetRoomByID(roomID)
	res.Room.RoomName = room.RoomName
	m.App.InfoLog.Println("ROOM:", res.Room.RoomName)
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w,r,"/make-reservation",http.StatusSeeOther)

}
