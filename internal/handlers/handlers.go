package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	//"log"
	// "io/ioutil"
	"io"
	"os"

	// "path/filepath"
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
	"golang.org/x/crypto/bcrypt"
//Charts
	"github.com/go-echarts/examples/examples"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

// Repo is repository used by handlers
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

// Testing function for repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// New Handlers sets the repository for handler
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
// (m *Repository) is the reciever function due to which all of the function is linked with AppConfig
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
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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

// Post Reservation handles the posting and validation of form
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

	htmlmessage := fmt.Sprintf(`
	<strong>Reservation Confirmation </strong><br>
	Dear %s,<br>
	This is to confirm your reservation from %s to %s
	
	`, reservation.FirstName, reservation.StartDate.Format("2006-01-01"), reservation.EndDate.Format("2006-01-02"))

	//send notification-first to  guest
	msg := models.MailData{
		To:       reservation.Email,
		From:     "Shobhitdimri7@gmail.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlmessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)
	m.App.Session.Put(r.Context(), "flash", "Email Sent!")
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// encrypts the password before saving it to database
func HashPassword(userPassword string) string {
	password, err := bcrypt.GenerateFromPassword([]byte(userPassword), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(password)
}

// Signing up new User
func (m *Repository) AdminSignup(w http.ResponseWriter, r *http.Request) {
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	accesslevel, _ := strconv.ParseInt(r.FormValue("access"), 10, 64)
	details := models.User{
		FirstName:   r.FormValue("first_name"),
		Email:       r.FormValue("email"),
		LastName:    r.FormValue("last_name"),
		Password:    r.FormValue("password"),
		AccessLevel: int(accesslevel),
	}
	fmt.Println("Access:-----------------------------", details.AccessLevel)
	if !form.Valid() {
		//TODO take user back to page
		m.App.Session.Put(r.Context(), "error", "Invalid!")
		render.Template(w, r, "admin_signup.html", &models.TemplateData{
			Form: form,
		})

		return
	}
	password := HashPassword(details.Password)
	details.Password = password

	AlreadyExist := m.DB.EmailCheck(details.Email)
	fmt.Println("status:", AlreadyExist)
	if AlreadyExist {
		m.App.Session.Put(r.Context(), "error", "User Exist")
		http.Redirect(w, r, "/admin/user", http.StatusSeeOther)
		return
	}
	err := m.DB.InsertUserData(details)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error")
		fmt.Println(err)
	}
	m.App.InfoLog.Println(details)
	m.App.Session.Put(r.Context(), "flash", "Signed up Successfully!")
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)

}

func (m *Repository) ShowAdmins(w http.ResponseWriter, r *http.Request) {
	users, err := m.DB.GetAllAdmins()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["users"] = users

	render.Template(w, r, "show_admin_user.html", &models.TemplateData{
		Data: data,
	})
}
func (m *Repository) ShowModifyAdminUsers(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	log.Println(id)
	users, err := m.DB.GetAdminByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["users"] = users
	log.Println(users)
	render.Template(w, r, "admin_edit.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}
func (m *Repository) AdminDeleteID(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	log.Println(id)
	err = m.DB.DeleteAdminByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Admin Deleted!!")
	http.Redirect(w, r, "/admin/data", http.StatusSeeOther)
}
func (m *Repository) UpdateAdminInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/admin/data", http.StatusSeeOther)
		return
	}
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	log.Println("------------------------",id)
	accesslevel, _ := strconv.ParseInt(r.FormValue("access"), 10, 64)
	fmt.Println("**************************************",accesslevel)
	details := models.User{
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Email:     r.FormValue("email"),
		AccessLevel: int(accesslevel),
		ID: id,
	}
	fmt.Println("showdata-------",details)
	err = m.DB.UpdateAdminData(details)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Admin Data Updated")
	http.Redirect(w, r, "/admin/data", http.StatusSeeOther)

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
	if(startDate==endDate){
		m.App.Session.Put(r.Context(), "warning", "Arrival & Departure Date Can't be Same")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
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
	m.App.InfoLog.Println("room_length---------------------------------------------------",len(rooms))
	
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

// AvailabilityJson handles request for availability and sends JSON as response
func (m *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	// need to parse request body
	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return appropriate json
		resp := jsonresponse{
			Ok:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		// got a database error, so return appropriate json
		resp := jsonresponse{
			Ok:      false,
			Message: "Error querying database",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	resp := jsonresponse{
		Ok:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	// I removed the error check, since we handle all aspects of
	// the json right here
	out, _ := json.MarshalIndent(resp, "", "     ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.html", &models.TemplateData{})
}

// Reservation summary displays the user reservation details
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

// Choose room displays lists tha available room to the user
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

// Bookroom takes url parameter build sessional variable and takes user to make reservation screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	//id , s , e
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")
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

	room, _ := m.DB.GetRoomByID(roomID)
	res.Room.RoomName = room.RoomName
	m.App.InfoLog.Println("ROOM:", res.Room.RoomName)
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "login.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// Handles user login Authentication
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Logging Working ...")
	//Renew Token prevent session fixation attack
	_ = m.App.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	fmt.Println(password)
	if !form.Valid() {
		//TODO take user back to page
		render.Template(w, r, "login.html", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, accesslevel, err := m.DB.Authenticate(email, password)

	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credential")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	render.Level = accesslevel
	// models.TemplateData.CurrentID := id
	intMap := make(map[string]int)
	intMap["active_id"] = id
	fmt.Println("Access_level-------", render.Level,"\nID:",intMap["active_id"])
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in Successfully")
	// new added
	render.Template(w, r, "home.html", &models.TemplateData{
		IntMap: intMap,
	})
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Destroys session
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	render.LoggedIn = false

	m.App.Session.Put(r.Context(), "warning", "Logged Out Successfully")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// Admin Dashboard Functions

var (
	itemCntPie = 2
	seasons    = []string{"Deluxe king Room", "Sunset Suite Room"}
)

func (m *Repository) CountRes() (int, int, int) {
	total_reservation_count, del_res, sun_res := m.DB.CountReservation()
	m.App.InfoLog.Println("count is :", total_reservation_count, "\nDeluxe King Room :", del_res, "\nSunset_Room : ", sun_res)
	return total_reservation_count, del_res, sun_res
}
func generatePieItems() []opts.PieData {
	_, b, c := Repo.CountRes()
	num := [2]int{b, c}
	items := make([]opts.PieData, 0)
	for i := 0; i < itemCntPie; i++ {
		items = append(items, opts.PieData{Name: seasons[i], Value: num[i]})
	}
	return items
}
func pieShowLabel() *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Reservation Count"}),
	)

	pie.AddSeries("pie", generatePieItems()).
		SetSeriesOptions(charts.WithLabelOpts(
			opts.Label{
				Show:      true,
				Formatter: "{b}: {c}",
			}),
		)
	return pie
}
//
var CurrentMonth []string 
var Data []int
func (m *Repository) MonthCount(){
	//to avoid duplicacy error in bar
	CurrentMonth = nil
	Data = nil
	Months := map[int]string{
     
		1: "Jan",
		2: "Feb",
		3: "Mar",
		4: "Apr",
		5: "May",
		6: "Jun",
		7: "Jul",
		8: "Aug",
		9: "Sep",
		10: "Oct",
		11: "Nov",
		12: "Dec",
	}

	month,BookingCount := m.DB.CountMonths()
	
	for i,_ := range month{
		CurrentMonth = append(CurrentMonth,Months[month[i]])
		m.App.InfoLog.Println("Mapping",Months[month[i]])
	}
	for i,_ := range BookingCount{
		Data = append(Data,BookingCount[i])
		m.App.InfoLog.Println("Data",BookingCount[i])
	}
} 
func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < len(CurrentMonth); i++ {
		items = append(items, opts.BarData{Value: Data[i]})
	}
	return items
}
func barWithTheme(theme string) *charts.Bar {
	year, _,_ := time.Now().Date()
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: theme}),
		charts.WithTitleOpts(opts.Title{Title: fmt.Sprintf("Room Count Month-Wise for Year - %d",year)}),
	)
	bar.SetXAxis(CurrentMonth).
		AddSeries("Category B", generateBarItems())
	return bar
}
func themeVintage() *charts.Bar {
	return barWithTheme(types.ThemeVintage)
}


func geoBase() *charts.Geo {
	geo := charts.NewGeo()
	geo.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "World Map"}),
		charts.WithGeoComponentOpts(opts.GeoComponent{
			Map:       "world",
			ItemStyle: &opts.ItemStyle{Color: "#006666"},
		}),
	)

	return geo
}

type PieExamples struct{}

func (PieExamples) Examples() {
	page := components.NewPage()
	page.AddCharts(
		pieShowLabel(),
		themeVintage(),
		geoBase(),
	)
	f, err := os.Create("templates/pie.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
func Users(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin_signup.html", &models.TemplateData{})
}
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	m.MonthCount()
	examplers := []examples.Exampler{
		PieExamples{},
	}
	for _, e := range examplers {
		e.Examples()
	}
	fmt.Println("Access_level -------", render.Level)
	render.Template(w, r, "admin.html", &models.TemplateData{})
	// render.Template(w, r, "admin_signup.html", &models.TemplateData{})
}


func (m *Repository) ChartLoad(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "pie.html", &models.TemplateData{})
}

// Shows all new reservation at Admin Page
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin_new_reservation.html", &models.TemplateData{
		Data: data,
	})
}

// Shows the entire Reservation database to Admin
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin_all_reservation.html", &models.TemplateData{
		Data: data,
	})
}

// AdminShow Reservation shows the Reservations at Admin Tool.
func (m *Repository) AdminShowReservations(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	log.Println(id)
	src := exploded[3]
	
	stringMap := make(map[string]string)
	stringMap["src"] = src
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap["month"] = month
	stringMap["year"] = year
	//get reservations from database....
	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "admin_show_reservations.html", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}
func (m *Repository) AdminPostShowReservations(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	log.Println(id)
	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src
	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")
	err = m.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	month := r.Form.Get("month")
	year := r.Form.Get("year")

	m.App.Session.Put(r.Context(), "flash", "Changes saved")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

// Marks a reservation as processed
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	_ = m.DB.UpdateProcessedForReservation(id, 1)
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")
	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// Delete Reservation from GUI
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	_ = m.DB.DeleteReservation(id)

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation deleted")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// Display Calendar to Admin
func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	// assume that there is no month/year specified
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]interface{})
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	// get the first and last days of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["rooms"] = rooms

	for _, x := range rooms {
		// create maps
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		// get all the restrictions for the current room
		restrictions, err := m.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		for _, y := range restrictions {
			if y.ReservationID > 0 {
				// it's a reservation
				for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}
			} else {
				// it's a block
				blockMap[y.StartDate.Format("2006-01-2")] = y.ID
			}
		}
		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap

		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)
	}

	render.Template(w, r, "admin_reservation_calender.html", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}

//Handle Reservation calender post
func (m *Repository)AdminPostReservationsCalendar(w http.ResponseWriter,r *http.Request){
	err := r.ParseForm()
	if err!=nil{
		helpers.ServerError(w,err)
		return
	}
	year,_ := strconv.Atoi(r.Form.Get("y"))
	month,_ := strconv.Atoi(r.Form.Get("m"))

	//processing Calendar form
	rooms,err := m.DB.AllRooms()
	if err!=nil{
		helpers.ServerError(w,err)
		return
	}
	form := forms.New(r.PostForm)
	for _,x := range rooms{
		//Get the block map for session loop through entire map & if we have entry in map 
		//that doesnot exist in posted map and if restriction id>0,then it is a block we need
		// to remove
		curMap := m.App.Session.Get(r.Context(),fmt.Sprintf("block_map_%d",x.ID)).(map[string]int)
		for name,value := range curMap{
			// var ok will be false if value is not in map
			if val,ok:=curMap[name];ok{
				//only pay attention to value > 0 and that are not in post form
				//the rest are placeholder for days w/o block
				if val>0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s",x.ID,name)){
						err := m.DB.DeleteBlockByID(value)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	}

	// now handle new blocks
	for name, _ := range r.PostForm {
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomID, _ := strconv.Atoi(exploded[2])
			t, _ := time.Parse("2006-01-2", exploded[3])
			// insert a new block
			err := m.DB.InsertBlockForRoom(roomID, t)
			if err != nil {
				log.Println(err)
			}
		}
	}

	m.App.Session.Put(r.Context(),"flash","Changes Saved")
	http.Redirect(w,r,fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d",year,month),http.StatusSeeOther)
}