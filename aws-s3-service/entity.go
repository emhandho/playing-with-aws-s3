package awss3service

import "time"

type ConfigSessionAWS struct {
	AwsURL             string
	AwsRegion          string
	AwsAccessKeyID     string
	AwsSecretAccessKey string
}

type BucketItems struct {
	Key          string
	LastModified time.Time
	Size         int64
	StorageClass string
}
