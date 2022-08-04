package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"
	"os"

	"github.com/Shobhitdimri01/Bookings/internal/config"
	"github.com/Shobhitdimri01/Bookings/internal/handlers"
	"github.com/Shobhitdimri01/Bookings/internal/models"
	"github.com/Shobhitdimri01/Bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)


const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	err := run()
	if err!=nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
func run() error{
		//Things I have stored in Session
		gob.Register(models.Reservation{})

		infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		app.InfoLog = infoLog

		errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
		app.ErrorLog = errorLog
		// change this to true when in production
		app.InProduction = false
	
		// set up the session
		session = scs.New()
		session.Lifetime = 24 * time.Hour
		session.Cookie.Persist = true
		session.Cookie.SameSite = http.SameSiteLaxMode
		session.Cookie.Secure = app.InProduction
	
		app.Session = session
	
		tc, err := render.CreateTemplateCache()
		if err != nil {
			log.Fatal("cannot create template cache")
			return err
		}
	
		app.TemplateCache = tc
		app.UseCache = false
	
		repo := handlers.NewRepo(&app)
		handlers.NewHandlers(repo)
	
		render.NewTemplates(&app)
		return nil
	}

