package awss3service

import "mime/multipart"

type Service interface {
	SaveConfig(input InputConfigAws) (ConfigSessionAWS, error)
	GetBucketsList() ([]string, error)
	ListBucketItems(bucketName string) ([]BucketItems, error)
	UploadFile(bucketName string, filename string, file multipart.File) error
	CreateBucket(bucketName string) error
	DeleteBucket(bucketName string) error
	DeleteItemInBucket(bucketName, itemName string) error
}