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

func AWSConfiguration(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join("views", "aws-config.html")
	var tmpl, err = template.ParseFiles(filepath, "views/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

    fmt.Println("Endpoint Hit: Set-Config Page")
}

func CreateBucketPage(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join("views", "create-bucket.html")
	var tmpl, err = template.ParseFiles(filepath, "views/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

    fmt.Println("Endpoint Hit: Create Bucket Page")
}