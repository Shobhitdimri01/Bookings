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
	
		mux.Get("/", handlers.Repo.Home)
		mux.Get("/about", handlers.Repo.About)
	
		return mux
	}
	