package handler

import (
	"aws-s3-sample/user"
	"fmt"
	"net/http"
)

type userHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) *userHandler {
	return &userHandler{service}
}

func (h *userHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// define struct to mapping the input from UI
	var input user.RegisterUser

	// secure method API
	if r.Method != "POST" {
		fmt.Println("Cannot Handle the kind of method type!")
	}

	// get the input value and map to struct
	input.Name = r.FormValue("name")
	input.Email = r.FormValue("email")
	input.Password = r.FormValue("password")

	_, err := h.service.RegisterUser(input)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Successfully Registering User!")
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}