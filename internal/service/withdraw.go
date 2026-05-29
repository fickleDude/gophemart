package service

import (
	"fmt"
	"strconv"

	"github.com/fickleDude/gophemart/internal/model"
	"github.com/fickleDude/gophemart/internal/repository"
)

type WithdrawService struct {
	repository *repository.WithdrawRepository
}

func NewWithdrawService(repository *repository.WithdrawRepository) *WithdrawService {
	return &WithdrawService{repository: repository}
}

// получение информации о выводе средств
func (w *WithdrawService) GetWithdraws(login string) ([]*model.Withdraw, error) {
	return w.repository.GetWithdraws(login)
}

// запрос на списание баллов + списание баллов
func (w *WithdrawService) AddWithdraw(withdraw model.Withdraw) error {
	return w.repository.AddWithdraw(withdraw.Login, withdraw.Order, withdraw.Sum)
}

// проверка номера на корректность (алгоритма Луна)
func (o *WithdrawService) ValidateOrder(number string) error {
	_, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return fmt.Errorf("неверный формат номера заказа")
	}
	return nil
}
