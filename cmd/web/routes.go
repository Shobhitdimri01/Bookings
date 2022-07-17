package main

import (
	"net/http"

	"github.com/Shobhitdimri01/Bookings/pkg/config"
	"github.com/Shobhitdimri01/Bookings/pkg/handlers"

	//"github.com/bmizerany/pat"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler  {
	//using pat router

	// mux :=pat.New()

	// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	// return mux

	//using chi router


		mux := chi.NewRouter()
	
		mux.Use(middleware.Recoverer)
		mux.Use(NoSurf)
		mux.Use(SessionLoad)

		//fileserver calls static template and helps to lad file(img)
		mux.Get("/", handlers.Repo.Home)
		mux.Get("/about", handlers.Repo.About)
		mux.Get("/deluxe-rooms", handlers.Repo.Generals)
		mux.Get("/suite-rooms", handlers.Repo.Majors)
		mux.Get("/search-availability", handlers.Repo.Availability)
		mux.Get("/contact", handlers.Repo.Contact)
	
		mux.Get("/make-reservation", handlers.Repo.Reservation)
	
		fileServer := http.FileServer(http.Dir("./static/"))
		mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	
		return mux
	}
	