package render

import (

	"net/http"
	"testing"

	"github.com/Shobhitdimri01/Bookings/internal/models"
)

func TestAddDefaultData(t *testing.T){
	var td models.TemplateData
	r , err  :=getSession()
	if err !=nil{
		t.Error(err.Error())
	}

	session.Put(r.Context(),"flash","123")
	result := AddDefaultData(&td,r)
	// if result.Flash == "123"{
	// 	t.Error("Mismatched Values")
	// }
	if result == nil{
		t.Error("Failed")
	}
}

func TestRenderTemplate(t *testing.T){
	pathtoTemplates ="./../../templates"
	tc,err := CreateTemplateCache()
	if err != nil {
		t.Error(err.Error())
	}
	app.TemplateCache = tc

	r,err := getSession()
	if err != nil{
		t.Error(err.Error())
	}

	var ww myWriter
	err = Template(&ww , r , "home.html",&models.TemplateData{})
	if err!=nil {
		t.Error("error writing templates to browser")
	}
	// err = RenderTemplate(&ww , r , "non-existing.html",&models.TemplateData{})
	// if err==nil {
	// 	t.Error("render the templates that doesn't exist")
	// }

}
func TestNewTemplates(t *testing.T){
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T)  {
	pathtoTemplates = "./../../templates"
	_,err:=CreateTemplateCache()
	if err!= nil{
		t.Error(err.Error())
	}
	
}

func getSession()(*http.Request,error){
	r,err := http.NewRequest("GET","/some-url",nil)
	if err != nil{
		return nil,err
	}
	ctx :=r.Context()
	ctx,_ = session.Load(ctx,r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r,nil
}