package repository

import (
	"database/sql"
	"errors"

	model "github.com/fickleDude/gophemart/internal/model"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// получение заказа по номеру
func (o *OrderRepository) GetOrder(number string) (*model.Order, error) {
	var order model.Order
	row := o.db.QueryRow("select number, uploaded_at, login from orders where number = $1", number)
	err := row.Scan(&order.Number, order.UploadedAt, order.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// получение списка загруженных пользователем номеров заказов
func (o *OrderRepository) GetOrders(login string) ([]*model.Order, error) {
	rows, err := o.db.Query(`select number, uploaded_at 
								from orders 
								where login = $1
								order by 2 desc`, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.Number, &order.UploadedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}
	return orders, nil
}

// загрузка пользователем номера заказа для расчёта
func (o *OrderRepository) AddOrder(login string, number string) error {
	_, err := o.db.Exec(`INSERT INTO orders (number, login, uploaded_at)
						 VALUES ($1, $2, now());`, number, login)
	if err != nil {
		return err
	}
	return nil
}
