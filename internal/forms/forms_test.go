package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T){
	r:= httptest.NewRequest("POST","/whatever",nil) //dummy url
	form :=New(r.PostForm)

	isValid := form.Valid()
	if isValid!=form.Valid(){
		t.Error("Error !! -  Got Invalid When Form should be Valid")
	}
}

func TestForm_Requires(t *testing.T){
	r:= httptest.NewRequest("POST","/whatever",nil)
	form :=New(r.PostForm)
	form.Required("A","b","c")
	if form.Valid(){
		t.Error("Error ! - Form shows Valid when required fields are missing")
	}
	postData := url.Values{}
	postData.Add("A","A")
	postData.Add("b","A")
	postData.Add("c","A")
	r,_ = http.NewRequest("POST","/whatever",nil)
	r.PostForm = postData
	form = New(r.PostForm)
	form.Required("A","b","c")
	if !form.Valid(){
		t.Error("Error ! - Got Invalid When Form should be Valid")
	}
}

func TestForm_Has(t *testing.T){
	
	postedData := url.Values{}
	form := New(postedData)
	x :=form.Has("a")
		if x{
			t.Error("form shows fields when it doesn't have")
		} 
		postedData = url.Values{}
		postedData.Add("a","a")
		form =  New(postedData)
	y := form.Has("a")
	if y != true{
		t.Error("Form doesn't show field data when it has One")
	}
	
}

func TestForm_Minlength(t *testing.T){
	postedData := url.Values{}
	form := New(postedData)
   form.Minlength("abcd",4)
   if form.Valid(){
	t.Error("Returning length of non-existance field")
   }
   //We have not added this data
		isError := form.Errors.Get("abcd")
		if isError == ""{
			t.Error("Data available yet showing null")
		}
   postedData = url.Values{}
		postedData.Add("add Data","Data Added") //[Key-> add Data; Value -> Data Added]
		form =  New(postedData)
		form.Minlength("add Data",100)
		if form.Valid(){
			t.Errorf("Form-data is valid for min 10 character but you passed %d as minimum Length",LengthContainer.Y)
		}


		postedData.Add("another field","abcdeer") //[Key-> add Data; Value -> Data Added]
		form =  New(postedData)
		form.Minlength("another field",1)
		if !form.Valid(){
			t.Errorf("Form-data is showing Incorrect min 1 character ! = passed %d as minimum Length",LengthContainer.Y)
		}
		//Data is added now
		isError = form.Errors.Get("another field")
		if isError != ""{
			t.Error("Data not available yet showing data")
		}
}

func TestForm_IsEmail(t *testing.T){
	postedData := url.Values{}
	form := New(postedData)
	form.IsEmail("x")
	if form.Valid(){
		t.Error("Form shows Valid when it doesn't have existing field")
	}
	postedData = url.Values{}
		postedData.Add("email","a@a.com") //[Key-> add Data; Value -> Data Added]
		form =  New(postedData)
		form.IsEmail("email")
		if !form.Valid(){
			t.Error("Got an Invalid Email When we should not have")
		}

		postedData = url.Values{}
		postedData.Add("invalid-email","a@invalid") //[Key-> add Data; Value -> Data Added]
		form =  New(postedData)
		form.IsEmail("invalid-email")
		if form.Valid(){
			t.Error("Got an Valid Email When we should not have")
		}
}
