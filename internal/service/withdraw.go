package service

import (
	"time"

	"github.com/fickleDude/gophemart/internal/helpers"
	"github.com/fickleDude/gophemart/internal/model"
	"github.com/fickleDude/gophemart/internal/repository"
)

type WithdrawServiceInterface interface {
	GetWithdraws(login string) ([]*model.Withdraw, error)
	AddWithdraw(withdraw model.Withdraw) error
	ValidateOrder(number string) bool
}

type WithdrawService struct {
	repository repository.WithdrawRepositoryInterface
}

func NewWithdrawService(repository repository.WithdrawRepositoryInterface) *WithdrawService {
	return &WithdrawService{repository: repository}
}

// получение информации о выводе средств
func (w *WithdrawService) GetWithdraws(login string) ([]*model.Withdraw, error) {
	return w.repository.GetWithdraws(login)
}

// запрос на списание баллов + списание баллов
func (w *WithdrawService) AddWithdraw(withdraw model.Withdraw) error {
	return w.repository.AddWithdraw(withdraw.Login, withdraw.Order, withdraw.Sum, time.Now().Format(time.RFC3339))
}

func (o *WithdrawService) ValidateOrder(number string) bool {
	return helpers.LuhnAlgorithm(number)
}
