package repository

import (
	"database/sql"
	"errors"

	model "github.com/fickleDude/gophemart/internal/model"
)

type InternalApiRepository struct {
	db *sql.DB
}

func NewInternalApiRepository(db *sql.DB) *InternalApiRepository {
	return &InternalApiRepository{db: db}
}

// получение информации о расчёте начислений баллов лояльности
func (o *InternalApiRepository) GetData(number string) (*model.Order, error) {
	var order model.Order
	row := o.db.QueryRow("select order_num, status, COALESCE(accrual, 0) from internal_service where order_num = $1", number)
	err := row.Scan(&order.Number, &order.Status, &order.Accrual)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}
