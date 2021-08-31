package repository

import (
	awsservice "aws-s3-sample/aws-s3-service"
	"database/sql"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(config awsservice.ConfigSessionAWS) error {
	// this is the for local storage aws s3 connection configuration
	// format := fmt.Sprintf("AWS_S3_URL=%s\nAWS_S3_REGION=%s\nAWS_S3_ACCESS_KEY_ID=%s\nAWS_S3_SECRET_ACCESS_KEY=%s\n",
	// 						config.AwsURL, config.AwsRegion, config.AwsAccessKeyID, config.AwsSecretAccessKey)
	// err := os.WriteFile("/Users/iCreativeLabs/Desktop/norman/icreativelabs-norman/playing-with-aws-s3/.env", []byte(format), 0644)
	// if err != nil {
	// 	return err
	// }

	// bellow using mysql for database
	stmt, err := r.db.Prepare("INSERT INTO users SET aws_url=?, aws_region=?, aws_access_key_id=?, aws_secret_access_key=?")
	if err == nil {
		_, err := stmt.Exec(&config.AwsURL, &config.AwsRegion, &config.AwsAccessKeyID, &config.AwsSecretAccessKey)
		if err != nil {
			return err
		}
	}

	return err
}