package config

import (
	"html/template"
	"log"

	"github.com/Shobhitdimri01/Bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

//holds the application config (struct) that we want to share with our entire application
type AppConfig struct {
	UseCache     			 bool
	TemplateCache 		map[string]*template.Template
	InProduction  			bool
	InfoLog						*log.Logger
	ErrorLog					*log.Logger
	Session     			  *scs.SessionManager
	MailChan 				   chan models.MailData
}
