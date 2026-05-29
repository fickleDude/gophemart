package service

import (
	"fmt"
	"strconv"

	model "github.com/fickleDude/gophemart/internal/model"
	"github.com/fickleDude/gophemart/internal/repository"
)

type OrderService struct {
	repository *repository.OrderRepository
}

func NewOrderService(repository *repository.OrderRepository) *OrderService {
	return &OrderService{repository: repository}
}

// получение заказа по номеру
func (o *OrderService) GetOrder(number string) (*model.Order, error) {
	return o.repository.GetOrder(number)
}

// получение списка загруженных пользователем номеров заказов
func (o *OrderService) GetOrders(login string) ([]*model.Order, error) {
	return o.repository.GetOrders(login)
}

// загрузка пользователем номера заказа для расчёта + пополнение баллов
func (o *OrderService) AddOrder(order model.Order) error {
	return o.repository.AddOrder(order.Login, order.Number)
}

// проверка номера на корректность (алгоритма Луна)
func (o *OrderService) ValidateOrder(number string) error {
	_, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return fmt.Errorf("неверный формат номера заказа")
	}
	return nil
}
