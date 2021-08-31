package user

type Service interface {
	RegisterUser(input RegisterUser) (User, error)
}