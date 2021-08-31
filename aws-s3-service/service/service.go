package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"

	awsservice "aws-s3-sample/aws-s3-service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	
)

type service struct {
	repository awsservice.Repository
}

func NewService(repo awsservice.Repository) *service {
	return &service{repo}
}

func (s *service) SaveConfig(input awsservice.InputConfigAws) (awsservice.ConfigSessionAWS ,error) {
	config := awsservice.ConfigSessionAWS{}
	config.AwsURL = input.AwsURL
	config.AwsRegion = input.AwsRegion
	config.AwsAccessKeyID = input.AwsAccessKeyID
	config.AwsSecretAccessKey = input.AwsSecretAccessKey

	err := s.repository.Save(config)
	if err != nil {
		return config, err
	}

	return config, nil
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
	// Create s3-service-client
	sess, err := s.createSession()
	if err != nil {
		return nil, errors.New("unable to create session")
	}
	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		return nil, fmt.Errorf("unable to list buckets, %v", err)
	}

	var bucketData []string
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
		return items, errors.New("unable to create session")
	}

	svc := s3.New(sess)

	response, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucketName)})
	if err != nil {
		return items, fmt.Errorf("unbale to list items in buckets %q, %v", bucketName, err)
	}

	if len(response.Contents) == 0 {
		return items, nil
	}

	for _, item := range response.Contents {
		itemt.Key = *item.Key
		itemt.LastModified = *item.LastModified
		itemt.Size = *item.Size
		itemt.StorageClass = *item.StorageClass
		itemt.BucketName = bucketName

		items = append(items, itemt)
	}

	return items, nil
}

func (s *service) CreateBucket(bucketName string) error {
	sess, err := s.createSession()
	if err != nil {
		return errors.New("unable to create session")
	}

	svc := s3.New(sess)
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("unable to create bucket %q, %v", bucketName, err)
	}

	fmt.Printf("Waiting for bucket %q to be created...\n", bucketName)

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("error occured while waiting for bucket to be created, %v", bucketName)
	}

	fmt.Printf("Bucket %q successfully created\n", bucketName)

	return nil
}

func (s *service) UploadFile(bucketName string, fileName string, fileLocation multipart.File) error {
	sess, err := s.createSession()
	if err != nil {
		return errors.New("unable to create session")
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   fileLocation,
	})
	if err != nil {
		return fmt.Errorf("unable to upload %q to %q, %v", fileName, bucketName, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", fileName, bucketName)
	return nil
}

func (s *service) DeleteBucket(bucketName string) error {
	sess, err := s.createSession()
	if err != nil {
		return errors.New("unable to create session")
	}

	// Create S3 service client
	svc := s3.New(sess)

	response, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucketName)})
	if err != nil {
		return fmt.Errorf("there is no items in this bucket %q, %v", bucketName, err)
	}

	if len((*response).Contents) != 0 {
		// make var that will be a place to store the objects
		objectsToDelete := make([]*s3.ObjectIdentifier, 0 , 1000)
		for _, object := range (*response).Contents{
			obj := s3.ObjectIdentifier {
				Key: object.Key,
			}
			objectsToDelete = append(objectsToDelete, &obj)
		}

		//Creating JSON payload for bulk delete
		deleteArray := s3.Delete{Objects: objectsToDelete}
		deleteParams := &s3.DeleteObjectsInput{
			Bucket: aws.String(bucketName),
			Delete: &deleteArray,
		}

		//Running the bulk delete job
		_, err := svc.DeleteObjects(deleteParams)
		if err != nil {
			return err
		}

		//Delete the bucket after empty
		_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("unable to delete bucket %q, %v", bucketName, err)
		}
	
		// Wait until bucket is deleted before finishing
		fmt.Printf("Waiting for bucket %q to be deleted...\n", bucketName)
		err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("error occurred while waiting for object %q to be deleted, %v", bucketName, err)
		}
	
		fmt.Printf("Bucket %q successfully deleted\n", bucketName)
		return nil
	}

	_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("unable to delete bucket %q, %v", bucketName, err)
	}

	// Wait until bucket is deleted before finishing
	fmt.Printf("Waiting for bucket %q to be deleted...\n", bucketName)

	err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return fmt.Errorf("error occurred while waiting for object %q to be deleted, %v", bucketName, err)
	}

	fmt.Printf("Bucket %q successfully deleted\n", bucketName)
	return nil
}

func (s *service) DeleteItemInBucket(bucketName, itemName string) error {
	sess, err := s.createSession()
	if err != nil {
		return errors.New("unable to create session")
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Delete the item
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(itemName),
	})
	if err != nil {
		return fmt.Errorf("unable to delete bucket %q, %v", itemName, err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(itemName),
	})
	if err != nil {
		return fmt.Errorf("error occurred while waiting for object %q to be deleted, %v", itemName, err)
	}

	fmt.Printf("Object %q successfully deleted\n", itemName)
	return nil
}
