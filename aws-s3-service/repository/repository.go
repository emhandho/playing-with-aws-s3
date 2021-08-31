package repository

import (
	awsservice "aws-s3-sample/aws-s3-service"
	"fmt"
	"os"
)

type repository struct {}

func NewRepository() *repository {
	return &repository{}
}

func (r *repository) Save(config awsservice.ConfigSessionAWS) error {
	format := fmt.Sprintf("AWS_S3_URL=%s\nAWS_S3_REGION=%s\nAWS_S3_ACCESS_KEY_ID=%s\nAWS_S3_SECRET_ACCESS_KEY=%s\n",
							config.AwsURL, config.AwsRegion, config.AwsAccessKeyID, config.AwsSecretAccessKey)
	err := os.WriteFile("/Users/iCreativeLabs/Desktop/norman/icreativelabs-norman/playing-with-aws-s3/.env", []byte(format), 0644)
	if err != nil {
		return err
	}

	return nil
}