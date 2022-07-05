package config

import (

	"html/template"
	"github.com/alexedwards/scs/v2"
)

//holds the application config (struct) that we want to share with our entire application
type AppConfig struct {
	UseCache     			 bool
	TemplateCache 		map[string]*template.Template
	InProduction  			bool
	Session     			  *scs.SessionManager
}
