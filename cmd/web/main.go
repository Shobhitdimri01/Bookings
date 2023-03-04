package main

import (
	// "database/sql/driver"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Shobhitdimri01/Bookings/internal/config"
	"github.com/Shobhitdimri01/Bookings/internal/driver"
	"github.com/Shobhitdimri01/Bookings/internal/handlers"
	"github.com/Shobhitdimri01/Bookings/internal/models"
	"github.com/Shobhitdimri01/Bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8082"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("Starting mail listener....")
	ListenforMail()

	fmt.Println("============================================================================================================")
	fmt.Printf("\t\tRunning application on port %s ...........\n", portNumber)
	fmt.Println("============================================================================================================")

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
func run() (*driver.DB, error) {
	//Things I have stored in Session
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})

	//creating channel for sending email
	mailchan := make(chan models.MailData)
	app.MailChan = mailchan

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

	//Connect to Database
	fmt.Println("\n\nConnecting to database.............")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=postgres")
	if err != nil {
		log.Fatal("Cannot connect to database ! dying .....", err.Error())
	}
	fmt.Println("=======================+++++++++++++++++++++++++++++++++++++++++++++========================================")
	fmt.Println("\t\t+----------------Succesfully Connected to PostgreSQL database----------------+")
	fmt.Println("=======================+++++++++++++++++++++++++++++++++++++++++++++========================================")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	return db, err
}
