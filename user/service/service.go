package service

import (
	"aws-s3-sample/user"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo user.Repository
}

func NewService(repo user.Repository) *service {
	return &service{repo}
}

func (s *service) RegisterUser(input user.RegisterUser) (user.User, error) {
	user := user.User{}
	user.Name = input.Name
	user.Email = input.Email

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordHash)

	newUser, err := s.repo.Save(user)
	if err != nil {
		return user, err
	}

	return newUser, nil
}
