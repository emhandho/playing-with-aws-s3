package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	repository "aws-s3-sample/aws-s3-service/repository"
	service "aws-s3-sample/aws-s3-service/service"
	handler "aws-s3-sample/handler"
	page "aws-s3-sample/handler/page-handler"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// func init() is the first func that will be executed when the program start
func init() {
	// use godot package to load/read the .env file and
	// return the value of the key
	// load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	handleRequests()
}

func handleRequests() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("username"), os.Getenv("password"), os.Getenv("address"), os.Getenv("dbname"))
	db, err := sql.Open(os.Getenv("dbdriver"), dsn)
	if err != nil {
		log.Fatal(err.Error())
	}

	
	awsRepo := repository.NewRepository(db)
	awsService := service.NewService(awsRepo)
	awsHandler := handler.NewConfigHandler(awsService)
	
	userRepo := repository.NewRepository(db)
	userService := service.NewService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	
	myRouter := mux.NewRouter().StrictSlash(true)
	
	myRouter.HandleFunc("/", page.HomePage)
	myRouter.HandleFunc("/login", page.LoginUser)
	myRouter.HandleFunc("/register", page.RegisterUser)
	myRouter.HandleFunc("/registeruser", userHandler.RegisterUser)
	myRouter.HandleFunc("/setconfig", page.AWSConfiguration)
	myRouter.HandleFunc("/createbucketpage", page.CreateBucketPage)
	myRouter.HandleFunc("/createbucket", awsHandler.CreateBucket)
	myRouter.HandleFunc("/uploadfile", awsHandler.UploadFile)
	myRouter.HandleFunc("/configaws", awsHandler.SetAWSConfiguration)
	myRouter.HandleFunc("/bucketslist", awsHandler.ListTheBuckets)
	myRouter.HandleFunc("/bucketdetails", awsHandler.ListBucketItems)
	myRouter.HandleFunc("/deletebucket", awsHandler.DeleteBucket)
	myRouter.HandleFunc("/deleteitem", awsHandler.DeleteItem)
	
	fmt.Println("Server running on port http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}