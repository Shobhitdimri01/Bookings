package handlers

import (
	"net/http"

	"github.com/Shobhitdimri01/Bookings/pkg/config"
	"github.com/Shobhitdimri01/Bookings/pkg/models"
	"github.com/Shobhitdimri01/Bookings/pkg/render"
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
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send data to the template
	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}