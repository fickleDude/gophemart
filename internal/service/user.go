package service

import (
	"crypto/sha256"
	"encoding/hex"

	model "github.com/fickleDude/gophemart/internal/model"
	"github.com/fickleDude/gophemart/internal/repository"
)

type UserServiceInterface interface {
	AddUser(user model.User) error
	GetUser(login string) (*model.User, error)
	ValidateUser(user model.User) (bool, error)
}

type UserService struct {
	repository repository.UserRepositoryInterface
}

func NewUserService(repository repository.UserRepositoryInterface) *UserService {
	return &UserService{repository: repository}
}

func getHash(password string) string {
	src := []byte(password)
	h := sha256.New()
	h.Write(src)
	dst := h.Sum(nil)
	return hex.EncodeToString(dst)
}

// регистрация пользователя
func (u *UserService) AddUser(user model.User) error {
	passwordHash := getHash(user.Password)
	return u.repository.AddUser(user.Login, passwordHash)
}

func (u *UserService) GetUser(login string) (*model.User, error) {
	return u.repository.GetUser(login)
}

// проверка пользователя
func (u *UserService) ValidateUser(user model.User) (bool, error) {
	repoUser, err := u.repository.GetUser(user.Login)
	if err != nil {
		return false, err
	}
	passwordHash := getHash(user.Password)
	isValid := repoUser.Password == passwordHash
	return isValid, nil
}
