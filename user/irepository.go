package user

type Repository interface {
	Save(User) (User, error)
}