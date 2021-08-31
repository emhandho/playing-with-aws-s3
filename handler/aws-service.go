package handler

import (
	awss3service "aws-s3-sample/aws-s3-service"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

type awsHandler struct {
	service awss3service.Service
}

func NewConfigHandler(service awss3service.Service) *awsHandler {
	return &awsHandler{service}
}

func (h *awsHandler) SetAWSConfiguration(w http.ResponseWriter, r *http.Request) {
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

	_, err := h.service.SaveConfig(input)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("successfully set the config!")

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (h *awsHandler) CreateBucket(w http.ResponseWriter, r *http.Request) {
	//get the bucket name from input
	inputName := r.FormValue("bucket_name")

	// input the name to the create bucket service
	err := h.service.CreateBucket(inputName)
	if err != nil {
		fmt.Println(errors.New("cannot create the bucket"))
	}

	fmt.Println("success create the bucket")

	http.Redirect(w, r, "/bucketslist", http.StatusMovedPermanently)
}

func (h *awsHandler) ListTheBuckets(w http.ResponseWriter, r *http.Request) {
	// get the buckets list from service
	buckets, err := h.service.GetBucketsList()
	if err != nil {
		fmt.Println(err.Error())
	}

	data := map[string]interface{}{
		"buckets": buckets,
	}

	tmpl, err := template.ParseFiles("views/buckets-list.html", "views/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Println("Endpoint Hit: buckets list page")
}

func (h *awsHandler) ListBucketItems(w http.ResponseWriter, r *http.Request) {
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
		"bucketName": bucketName,
		"data":       items,
		"message":    message,
		"noData":     noData,
	}

	tmpl, err := template.ParseFiles("views/bucket-item.html", "views/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Println("Endpoint Hit: bucket items list page")
}

func (h *awsHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	content := r.Header.Get("Content-Type")
	fmt.Println(content)

	// get the params from input
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("here")
		fmt.Println(err)
		return
	}

	bucketName := r.FormValue("name")
	uploadedFile, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	fileName := handler.Filename

	if err := h.service.UploadFile(bucketName, fileName, uploadedFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Successfully Upload file %s to aws bucket %s\n", fileName, bucketName)
	http.Redirect(w, r, "/bucketdetails?name="+bucketName, http.StatusMovedPermanently)
}

func (h *awsHandler) DeleteBucket(w http.ResponseWriter, r *http.Request) {
	// get the bucket name
	bucketName := r.FormValue("name")

	// input the name to the deelet bucket service
	err := h.service.DeleteBucket(bucketName)
	if err != nil {
		fmt.Println(errors.New("cannot delete bucket"))
	}

	fmt.Println("success create the bucket")

	http.Redirect(w, r, "/bucketslist", http.StatusMovedPermanently)
}

func (h *awsHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	// get the bucket item name
	itemName := r.FormValue("item")
	bucketName := r.FormValue("name")

	if err := h.service.DeleteItemInBucket(bucketName, itemName); err != nil {
		fmt.Println(errors.New("cannot delete item"))
	}

	fmt.Println("success delete the item")

	http.Redirect(w, r, "/bucketdetails?name="+bucketName, http.StatusMovedPermanently)
}
