package service

import (
	model "github.com/fickleDude/gophemart/internal/model"
	"github.com/fickleDude/gophemart/internal/repository"
)

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{repository: repository}
}

// регистрация пользователя
func (u *UserService) AddUser(user model.User) error {
	return u.repository.AddUser(user.Login, user.Password)
}

// проверка пользователя
func (u *UserService) GetUser(login string) (*model.User, error) {
	return u.repository.GetUser(login)
}
