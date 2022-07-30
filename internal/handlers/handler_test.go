package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	Value string
}

//TheTest is testing format for our routes
var theTests = []struct {
	name               string
	url                string
	method             string
	param              []postData
	expectedStatusCode int
}{
	{
		"home",    //Homepage Route
		"/",		//Url
		"GET",		//Get Request
		[]postData{},   //Home require No parameter
		http.StatusOK, //200
	},
	//Testing all the Get pages 
	{"About", "/about", "GET", []postData{}, http.StatusOK},
	{"Deluxe", "/deluxe-rooms", "GET", []postData{}, http.StatusOK},
	{"Suite", "/suite-rooms", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"make-res", "/make-reservation", "GET", []postData{}, http.StatusOK},


}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {

		}
	}
}
