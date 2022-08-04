package render

import (
	"github.com/Shobhitdimri01/Bookings/internal/config"
	"github.com/alexedwards/scs/v2"
	"os"
	"testing"
	"encoding/gob"
	"github.com/Shobhitdimri01/Bookings/internal/models"
	"time"
	"net/http"
	"log"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M){
	gob.Register(models.Reservation{})

		// change this to true when in production
		testApp.InProduction = false

		infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		testApp.InfoLog = infoLog

		errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
		testApp.ErrorLog = errorLog
		// set up the session
		session = scs.New()
		session.Lifetime = 24 * time.Hour
		session.Cookie.Persist = true
		session.Cookie.SameSite = http.SameSiteLaxMode
		session.Cookie.Secure = false
	
		testApp.Session = session

		app = &testApp


	os.Exit(m.Run())//Just before the application closes it runs our test
}

//Creating our response Writter struct

type myWriter struct{

}

func (tw *myWriter)Header() http.Header{
	var h http.Header
	return h
}
func(tw *myWriter) WriteHeader( i int){

}
func (tw *myWriter) Write(b []byte)(int , error){
	length := len(b)
	return length,nil
}