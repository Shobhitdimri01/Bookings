package main

import (
	"fmt"
	"testing"

	"github.com/Shobhitdimri01/Bookings/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux:=routes(&app)

	switch v := mux.(type){
	case *chi.Mux:
		//do nothing Testcase passed
	default:
		t.Error(fmt.Sprintf("type is not *chi.Mux , but of type %T",v))
	}
}