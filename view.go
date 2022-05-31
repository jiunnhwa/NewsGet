package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

var tplDir = "./html/templates"

//ViewData is a collection of data for the view
type ViewData struct {
	RowCount            int
	Message             string
	PageTitle           string
	ResponseTitle       string
	ResponseBody        string
	ResponseDescription string
	HasSessionID        bool
	IsLoggedIn          bool
	IsActionCreateJob   bool

	DataURL string
	URL     string

	Name string `json:"name"`
	Msg  string `json:"msg"`

	LogLines        string
	WebSocketOutput string

	DynamicList map[string]interface{}
	DynamicMap  map[string]interface{}
}

type View struct {
	Files []string
	tpl   *template.Template
	ViewData
}

//NewView constructs the view with parsing the files
func NewView(files []string) *View {
	v := &View{Files: files}
	v.ParseFiles()
	return v
}

//ParseFiles and returns View
func (v *View) ParseFiles() *View {
	tmpl, err := template.ParseFiles(v.Files...)
	v.tpl = template.Must(tmpl, err)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return v
}

//SetViewData updates ViewData
func (v *View) SetViewData(vd *ViewData) *View {
	v.ViewData = *vd
	return v
}

//ServeTemplate writes out from template
func (v *View) ServeTemplate(w http.ResponseWriter, r *http.Request) *View {
	v.tpl.Execute(w, v.ViewData)
	return v
}

//ToHTML writes out the template as html
func (v *View) ToHTML(ffname string) *View {
	file, err := os.Create(ffname)
	if err != nil {
		return v
	}
	defer file.Close()
	v.tpl.Execute(file, v.ViewData)
	return v
}
