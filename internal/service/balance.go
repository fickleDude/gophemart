package service

import (
	"net/http"

	"github.com/fickleDude/gophemart/internal/helpers"
	model "github.com/fickleDude/gophemart/internal/model"
	"github.com/fickleDude/gophemart/internal/repository"
)

type BalanceServiceInterface interface {
	GetBalance(login string) (*model.Balance, error)
}

type BalanceService struct {
	orderRepository    *repository.OrderRepository
	withdrawRepository *repository.WithdrawRepository
}

func NewBalaneService(orderRepository *repository.OrderRepository, withdrawRepository *repository.WithdrawRepository) *BalanceService {
	return &BalanceService{orderRepository: orderRepository, withdrawRepository: withdrawRepository}
}

// получение текущего баланса
func (b *BalanceService) GetBalance(login string) (*model.Balance, error) {
	//get all orders
	orders, err := b.orderRepository.GetOrders(login)
	if err != nil {
		return nil, err
	}
	//get all withdraws
	withdraws, err := b.withdrawRepository.GetWithdraws(login)
	if err != nil {
		return nil, err
	}

	//count accrual
	accrual := 0.0
	client := http.Client{}
	for _, o := range orders {
		orderDetails, err := helpers.GetOrderAccrual(o.Number, client)
		if err != nil {
			return nil, err
		}
		if orderDetails.Status == "PROCESSED" {
			accrual += orderDetails.Accrual
		}
	}
	//count withdraws
	withdraw := 0.0
	for _, w := range withdraws {
		withdraw += w.Sum
	}

	balance := model.Balance{Current: accrual - withdraw, Withdraw: withdraw}
	return &balance, nil
}
