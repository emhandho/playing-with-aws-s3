package awss3service

type Repository interface {
	Save(Config ConfigSessionAWS) error
}