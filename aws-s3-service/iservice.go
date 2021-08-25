package awss3service

type Service interface {
	SaveConfig(input InputConfigAws) error
	GetBucketsList() ([]string, error)
	ListBucketItems(bucketName string) ([]BucketItems, error)
}