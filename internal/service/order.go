package service

import (
	"net/http"
	"time"

	"github.com/fickleDude/gophemart/internal/api"
	"github.com/fickleDude/gophemart/internal/helpers"
	model "github.com/fickleDude/gophemart/internal/model"
	"github.com/fickleDude/gophemart/internal/repository"
)

type OrderServiceInterface interface {
	GetOrder(number string) (*model.Order, error)
	GetOrders(login string) ([]*model.Order, error)
	AddOrder(order model.Order) error
	ValidateOrder(number string) bool
}

type OrderService struct {
	repository repository.OrderRepositoryInterface
}

func NewOrderService(repository repository.OrderRepositoryInterface) *OrderService {
	return &OrderService{repository: repository}
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
		orderDetails, err := api.GetOrderAccrual(o.Number, client)
		if err != nil {
			return nil, err
		}
		if orderDetails != nil {
			o.Accrual = orderDetails.Accrual
			o.Status = orderDetails.Status
		}
	}
	return orders, nil
}

// загрузка пользователем номера заказа для расчёта + пополнение баллов
func (o *OrderService) AddOrder(order model.Order) error {
	return o.repository.AddOrder(order.Login, order.Number, time.Now().Format(time.RFC3339))
}

func (o *OrderService) ValidateOrder(number string) bool {
	return helpers.LuhnAlgorithm(number)
}
