package awss3service

import "time"

type ConfigSessionAWS struct {
	AwsURL             string
	AwsRegion          string
	AwsAccessKeyID     string
	AwsSecretAccessKey string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type BucketItems struct {
	BucketName   string
	Key          string
	LastModified time.Time
	Size         int64
	StorageClass string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
