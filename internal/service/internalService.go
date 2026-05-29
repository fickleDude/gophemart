package service

import (
	model "github.com/fickleDude/gophemart/internal/model"
	"github.com/fickleDude/gophemart/internal/repository"
)

type InternalApiService struct {
	repository *repository.InternalApiRepository
}

func NewInternalApiService(repository *repository.InternalApiRepository) *InternalApiService {
	return &InternalApiService{repository: repository}
}

// регистрация пользователя
func (i *InternalApiService) GetData(number string) (*model.Order, error) {
	return i.repository.GetData(number)
}
