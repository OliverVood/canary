package template

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type Data struct {
	Title string
	Data  any
}

func NewData(data any, title string) *Data {
	return &Data{
		Title: title,
		Data:  data,
	}
}

func Render(w http.ResponseWriter, name string, data any) {
	layout := filepath.Join(".", "cmd", "internal", "ui", "views", "layouts", "index.html")
	view := filepath.Join(".", "cmd", "internal", "ui", "views", fmt.Sprintf("%s.html", name))

	tmpl, err := template.ParseFiles(layout, view)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, "Ошибка рендеринга: "+err.Error(), http.StatusInternalServerError)
	}
}
