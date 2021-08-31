package repository

import (
	"aws-s3-sample/user"
	"database/sql"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user user.User) (user.User, error) {
	// using mysql, prepared statement for avoiding from SQL injection
	stmt, err := r.db.Prepare("INSERT INTO users SET username=?, password=?, fullname=?, email=?")
	if err == nil {
		_, err := stmt.Exec(&user.Name, &user.Email, &user.PasswordHash)
		if err != nil {
			return user, err
		}
	}
	return user, err
}
