package main

import (
	"fmt"
	"log"
	"net/http"

	repository "aws-s3-sample/aws-s3-service/repository"
	service "aws-s3-sample/aws-s3-service/service"
	page "aws-s3-sample/handler/page-handler"
	handler "aws-s3-sample/handler"

	"github.com/gorilla/mux"
)

func main() {
	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	awsRepo := repository.NewRepository()
	awsService := service.NewService(awsRepo)
	awsHandler := handler.NewConfigHandler(awsService)

	myRouter.HandleFunc("/", page.HomePage)
	myRouter.HandleFunc("/configaws", awsHandler.SetAWSConfiguration)
	myRouter.HandleFunc("/bucketslist", awsHandler.ListTheBuckets)
	myRouter.HandleFunc("/bucketdetails", awsHandler.ListBucketItems)
	
	fmt.Println("Server running on port http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}