package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"log"

	"github.com/Shobhitdimri01/Bookings/internal/config"
	"github.com/Shobhitdimri01/Bookings/internal/models"
	"github.com/Shobhitdimri01/Bookings/internal/render"
	"github.com/Shobhitdimri01/Bookings/internal/forms"
)

//Repo is repository used by handlers
var Repo *Repository

//Repository is a type of repo
type Repository struct{
	App  *config.AppConfig
}



//New Repo creates the new repository
func NewRepo(a *config.AppConfig) *Repository{
	return &Repository{
		App: a,
	}
}


//New Handlers sets the repository for handler
func NewHandlers(r *Repository){
	Repo = r
}

// Home is the handler for the home page
//(m *Repository) is the reciever function due to which all of the function is linked with AppConfig
// Home is the handler for the home page
// Home is the handler for the home page
// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.html", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send data to the template
	render.RenderTemplate(w, r, "about.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r , "make-reservation.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}
//Post Reservation handles the posting and validation of form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err :=r.ParseForm()
	if err != nil {
		log.Println(err.Error())
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName: r.Form.Get("last_name"),
		Email: r.Form.Get("email"),
		Phone: r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)
	//form.Has("first_name", r)
	form.Required("first_name","last_name","email","phone")
	form.Minlength("first_name",3,r)
	form.IsEmail("email")

	if !form.Valid(){
		data := make(map[string] interface{})
		data ["reservation"] = reservation

		render.RenderTemplate(w, r , "make-reservation.html", &models.TemplateData{
			Form : form,
			Data:   data, 
		})
	}

	//Calling Session to store Reservation details
	m.App.Session.Put(r.Context(),"reservation",reservation)

	http.Redirect(w,r,"/summary",http.StatusSeeOther)

}

// Generals renders the room page
func (m *Repository) Deluxe(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r,  "deluxe.html", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Suite(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r,  "suite.html", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r ,  "search-availability.html", &models.TemplateData{})
}

// PostAvailability renders the search availability page
// PostAvailability handles post
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end is %s", start, end)))
}



type jsonresponse struct{
	Ok			 bool 		`json:"ok"`
	Message string		`json:"message"`
}
//AvailabilityJson handles request for availability and sends JSON as response
func (m *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	resp := jsonresponse{
			Ok : 				true,
			Message: 	"Available!",
	}
	out , err := json.MarshalIndent(resp,"","    ")
	if err != nil {
		log.Println(err.Error())
	}

	// log.Println(string(out))
	w.Header().Set("Content-Type","application/json")
	w.Write(out)

}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r,  "contact.html", &models.TemplateData{})
}

//Reservation summary
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	//Pulling the data out from Session
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("can't get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(w, r, "reservation-summary.html", &models.TemplateData{
		Data: data,
	})
}

