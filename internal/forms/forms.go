package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

//Creates a custom form struct,embeds a url Values object
type Form struct {
	url.Values
	Errors 		errors
}

//Checking forms Validation
//Valid return true if there are no errors , otherwise return false
func (f *Form)Valid() bool {
	return len(f.Errors)==0
}

//New initializes a form struct
func New(data url.Values) *Form{
	return &Form{
		data,
		errors(map [string][]string{}),
	}
}

//This is function that checks whether all the fields are filled by user in Make-Reservation Page.
func(f *Form) Required (fields ...string){
	for _,field := range fields{
		value :=f.Get(field)
		if strings.TrimSpace(value) == ""{
			f.Errors.Add(field,"This field can't be left blank")
		}
	}
}

//Checks if required form field is in post and not empty
func (f *Form) Has (field string, r *http.Request)bool{
	 x := r.Form.Get(field)
	 if x =="" {
		
		return false
	 }
	 return true
}

//Checks minimum length
func (f *Form) Minlength (field string, length int, r *http.Request ) bool {
	x := r.Form.Get(field)
	if len(x) <length {
		f.Errors.Add(field , fmt.Sprintf("This field should be atleast %d character long",length))
		return false
	}
	return true
}

func (f *Form)IsEmail(field string){
	if !govalidator.IsEmail(f.Get(field)){
		f.Errors.Add(field,"Invalid Email !!!")
	}
}