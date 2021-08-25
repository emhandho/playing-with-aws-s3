package pagehandler

import (
	"net/http"
	"html/template"
	"path"
	"fmt"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join("views", "home.html")
	var tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

    fmt.Println("Endpoint Hit: Home Page")
}