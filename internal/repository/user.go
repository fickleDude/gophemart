package repository

import (
	model "github.com/fickleDude/gophemart/internal/model"
)

type UserRepository struct {
	storage []*model.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{storage: []*model.User{
		{Login: "123", Password: "9278923470"},
	}}
}

// регистрация пользователя
func (u *UserRepository) AddUser(login string, password string) error {
	user := model.User{Login: login, Password: password}
	u.storage = append(u.storage, &user)
	return nil
}
