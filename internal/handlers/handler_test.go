package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	//Testing all the Get pages 
	{
		"home",    //Homepage Route
		"/",		//Url
		"GET",		//Get Request
		[]postData{},   //Home require No parameter
		http.StatusOK, //200
	},
	{"About", "/about", "GET", []postData{}, http.StatusOK},
	{"Deluxe", "/deluxe-rooms", "GET", []postData{}, http.StatusOK},
	{"Suite", "/suite-rooms", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"make-res", "/make-reservation", "GET", []postData{}, http.StatusOK},
	
	//Testing all the Post routes
	{"post-search-availability", "/search-availability", "POST", []postData{
		{key: "start" , Value: "12-10-2021"},
		{key:"end", Value:"13-10-2021"},

	}, http.StatusOK},

	{"post-search-availability-json", "/search-availability-json", "POST", []postData{
		{key: "start" , Value: "12-10-2021"},
		{key:"end", Value:"13-10-2021"},

	}, http.StatusOK},
	{"post-reservation", "/make-reservation", "POST", []postData{
		{key: "first_name" , Value: "Rahul"},
		{key:"last_name", Value:"Dewan"},
		{key:"email", Value:"Rahul@gmail.com"},
		{key:"phone", Value:"9192921919"},

	}, http.StatusOK},

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
		} else {// e.Method == "post"
				values :=url.Values{}
				//ranging over the param value in []postdata
				for _,x := range e.param {
					values.Add(x.key , x.Value)
				}
				resp , err :=ts.Client().PostForm(ts.URL+e.url, values)
				if err!=nil{
					t.Log(err.Error())
					t.Fatal(err.Error())
				} 
				if resp.StatusCode != e.expectedStatusCode {
					t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
				}
		}
	}
}
