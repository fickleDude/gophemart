package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	model "github.com/fickleDude/gophemart/internal/model"
	"github.com/fickleDude/gophemart/internal/repository"
)

type OrderServiceInterface interface {
	GetOrder(number string) (*model.Order, error)
	GetOrders(login string) ([]*model.Order, error)
	AddOrder(order model.Order) error
	ValidateOrder(number string) error
}

type OrderService struct {
	repository *repository.OrderRepository
}

func NewOrderService(repository *repository.OrderRepository) *OrderService {
	return &OrderService{repository: repository}
}

func getOrderAccrual(number string, client http.Client) (*model.Order, error) {
	baseURL := "http://localhost:8081/api/orders"
	url := fmt.Sprintf("%s/%s", baseURL, number)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	var order model.Order
	if err := json.NewDecoder(response.Body).Decode(&order); err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return &order, nil

}

// получение заказа по номеру
func (o *OrderService) GetOrder(number string) (*model.Order, error) {
	return o.repository.GetOrder(number)
}

// получение списка загруженных пользователем номеров заказов
func (o *OrderService) GetOrders(login string) ([]*model.Order, error) {
	orders, err := o.repository.GetOrders(login)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	for _, o := range orders {
		orderDetails, err := getOrderAccrual(o.Number, client)
		if err != nil {
			return nil, err
		}
		o.Accrual = orderDetails.Accrual
		o.Status = orderDetails.Status
	}
	return orders, nil
}

// загрузка пользователем номера заказа для расчёта + пополнение баллов
func (o *OrderService) AddOrder(order model.Order) error {
	return o.repository.AddOrder(order.Login, order.Number, time.Now().Format(time.RFC3339))
}

// проверка номера на корректность (алгоритма Луна)
func (o *OrderService) ValidateOrder(number string) error {
	_, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return fmt.Errorf("неверный формат номера заказа")
	}
	return nil
}
