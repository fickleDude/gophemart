package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/fickleDude/gophemart/internal/model"
	"github.com/fickleDude/gophemart/internal/repository"
)

type BalanceService struct {
	orderRepository    *repository.OrderRepository
	withdrawRepository *repository.WithdrawRepository
}

func NewBalaneService(orderRepository *repository.OrderRepository, withdrawRepository *repository.WithdrawRepository) *BalanceService {
	return &BalanceService{orderRepository: orderRepository, withdrawRepository: withdrawRepository}
}

func getOrderAccrual(number string, client http.Client) (float64, error) {
	baseURL := "http://localhost:8081/api/orders"
	url := fmt.Sprintf("%s/%s", baseURL, number)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	var order model.Order
	if err := json.NewDecoder(response.Body).Decode(&order); err != nil {
		return 0, err
	}
	defer response.Body.Close()
	if order.Status == "PROCESSED" {
		return order.Accrual, nil
	}
	return 0, nil

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
		delta, err := getOrderAccrual(o.Number, client)
		if err != nil {
			return nil, err
		}
		accrual += delta
	}
	//count withdraws
	withdraw := 0.0
	for _, w := range withdraws {
		withdraw += w.Sum
	}

	balance := model.Balance{Current: accrual - withdraw, Withdraw: withdraw}
	return &balance, nil
}
