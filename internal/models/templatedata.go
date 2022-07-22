package models

//Data that will be used and send to Template pages
//Template Data holds data send from handlers to templates
type TemplateData struct{
	StringMap  		   map[string]string
	IntMap		  		map[string]int
	FloatMap    	   map[string]float32
	Data 				  map[string]interface{}
	CSRFToken       string
	Flash				  string
	Warning			   string
	Error                string
}