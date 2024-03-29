package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Shobhitdimri01/Bookings/internal/config"
	"github.com/Shobhitdimri01/Bookings/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"humanDate": HumanDate,
	"formatDate":FormatDate,
	"iterate":iterate,
	"add":Add,
	"userid":GetuserID,
}
var pathtoTemplates = "./templates"

var app *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	app = a
}

func Add(a,b int)int{
	return a+b
}
//Iterates return slice of ints
func iterate(count int)[]int{
	var(
		i int
	    items []int
	)
	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}
func HumanDate(t time.Time) string {
	return t.Format("02-01-2006")
}

func FormatDate(t time.Time,f string)string{
	return t.Format(f)
}
var Userid string
func GetuserID()string{
	return Userid
}
var( LoggedIn bool 
	Level int)

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	td.Level = Level
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
		LoggedIn = true
	}
	// app.InfoLog.Println("LoggedIn", LoggedIn)
	return td
}

// td is template data
// Template renders a template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
		return errors.New("Can't get the templates")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}
	return nil
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.html", pathtoTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathtoTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathtoTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
