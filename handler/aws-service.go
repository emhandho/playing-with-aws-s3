package handler

import (
	awss3service "aws-s3-sample/aws-s3-service"
	"fmt"
	"html/template"
	"net/http"
)

type AWSHandler struct {
	service awss3service.Service
}

func NewConfigHandler(service awss3service.Service) *AWSHandler {
	return &AWSHandler{service}
}

func (h *AWSHandler) SetAWSConfiguration(w http.ResponseWriter, r *http.Request) {
	// define the struct to mapping the input
	var input awss3service.InputConfigAws

	// check if the method is POST or not
	if r.Method != "POST" {
		fmt.Println("cannot handle the method type!")
	}

	// get the input value
	awsURL := r.FormValue("aws_url")
	awsRegion := r.FormValue("aws_region")
	awsAccessKeyID := r.FormValue("aws_access_key_id")
	awsSecretAccessKey := r.FormValue("aws_secret_access_key")

	input.AwsURL = awsURL
	input.AwsRegion = awsRegion
	input.AwsAccessKeyID = awsAccessKeyID
	input.AwsSecretAccessKey = awsSecretAccessKey

	err := h.service.SaveConfig(input)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("successfully set the config!")

	http.Redirect(w, r, "/bucketslist", http.StatusMovedPermanently)
}

func (h *AWSHandler) ListTheBuckets(w http.ResponseWriter, r *http.Request) {
	// get the buckets list from service
	buckets, err := h.service.GetBucketsList()
	if err != nil {
		fmt.Println(err.Error())
	}

	data := map[string]interface{}{
		"buckets": buckets,
	}

	tmpl, err := template.ParseFiles("views/bucketslist.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Println("Endpoint Hit: buckets list page")
}

func (h *AWSHandler) ListBucketItems(w http.ResponseWriter, r *http.Request) {
	// get the params for service
	bucketName := r.FormValue("name")

	// call the service to get bucket details
	items, err := h.service.ListBucketItems(bucketName)
	if err != nil {
		fmt.Println(err.Error())
	}

	var message string
	var noData bool
	if len(items) == 0 {
		message = "There is no items in this bucket."
		noData = false
	} else {
		message = ""
		noData = true
	}

	data := map[string]interface{}{
		"data": items,
		"message" : message,
		"noData" : noData,
	}

	tmpl, err := template.ParseFiles("views/bucketitem.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Println("Endpoint Hit: bucket items list page")
}
