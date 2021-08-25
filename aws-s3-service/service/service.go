package service

import (
	"fmt"
	"log"
	"os"

	awsservice "aws-s3-sample/aws-s3-service"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	// "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type service struct {
	repository awsservice.Repository
}

func NewService(repo awsservice.Repository) *service {
	return &service{repo}
}

func getConfig() {
	// use godot package to load/read the .env file and
	// return the value of the key
	// load .env file
	err := godotenv.Load("save-data/config.env")
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalf("Error loading .env file")
	}
}

func (s *service) SaveConfig(input awsservice.InputConfigAws) error {
	config := awsservice.ConfigSessionAWS{}
	config.AwsURL = input.AwsURL
	config.AwsRegion = input.AwsRegion
	config.AwsAccessKeyID = input.AwsAccessKeyID
	config.AwsSecretAccessKey = input.AwsSecretAccessKey

	err := s.repository.Save(config)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) createSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:                        aws.String(os.Getenv("AWS_S3_REGION")),
		Endpoint:                      aws.String(os.Getenv("AWS_S3_URL")),
		CredentialsChainVerboseErrors: aws.Bool(true),
		S3ForcePathStyle:              aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_S3_ACCESS_KEY_ID"),
			os.Getenv("AWS_S3_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return sess, nil
}

func (s *service) GetBucketsList() ([]string, error) {
	// get the config first
	getConfig()

	// Create s3-service-client
	sess, err := s.createSession()
	if err != nil {
		s.exitErrorf("Unable to create session")
	}
	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		s.exitErrorf("Unable to list buckets, %v", err)
	}

	var bucketData []string
	fmt.Println("Buckets Name:")
	for _, bucket := range result.Buckets {
		bucketData = append(bucketData, aws.StringValue(bucket.Name))
	}
	
	return bucketData, nil
}

func (s *service) ListBucketItems(bucketName string) ([]awsservice.BucketItems, error) {
	// define struct for input
	var items []awsservice.BucketItems
	var itemt awsservice.BucketItems

	sess, err := s.createSession()
	if err != nil {
		s.exitErrorf("Unable to create session")
	}

	svc := s3.New(sess)

	response, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucketName)})
	if err != nil {
		s.exitErrorf("Unbale to list items in buckets %q, %v", bucketName, err)
	}

	if len(response.Contents) > 1 {
		for _, item := range response.Contents {
			itemt.Key = *item.Key
			itemt.LastModified = *item.LastModified
			itemt.Size = *item.Size
			itemt.StorageClass = *item.StorageClass
			
			items = append(items, itemt)
		}
	} else {
		for _, item := range response.Contents {
			itemt.Key = *item.Key
			itemt.LastModified = *item.LastModified
			itemt.Size = *item.Size
			itemt.StorageClass = *item.StorageClass
		
			items = append(items, itemt)
		}
	}

	return items, nil
}

func (s *service) exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
