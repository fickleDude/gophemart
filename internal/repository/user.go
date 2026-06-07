package repository

import (
	"database/sql"
	"errors"

	model "github.com/fickleDude/gophemart/internal/model"
)

type UserRepositoryInterface interface {
	AddUser(login string, password string) error
	GetUser(login string) (*model.User, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// регистрация пользователя
func (u *UserRepository) AddUser(login string, password string) error {
	_, err := u.db.Exec(`INSERT INTO users (login, password)
						 VALUES ($1, $2);`, login, password)
	if err != nil {
		return err
	}
	return nil
}

// проверка пользователя
func (u *UserRepository) GetUser(login string) (*model.User, error) {
	var user model.User
	row := u.db.QueryRow("select login, password from users where login = $1", login)
	err := row.Scan(&user.Login, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
