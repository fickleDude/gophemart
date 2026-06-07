package repository

import (
	"database/sql"

	"github.com/fickleDude/gophemart/internal/model"
)

type WithdrawRepositoryInterface interface {
	GetWithdraws(login string) ([]*model.Withdraw, error)
	AddWithdraw(login string, order string, sum float64, processedAt string) error
}

type WithdrawRepository struct {
	db *sql.DB
}

func NewWithdrawRepository(db *sql.DB) *WithdrawRepository {
	return &WithdrawRepository{db: db}
}

// получение информации о выводе средств
func (w *WithdrawRepository) GetWithdraws(login string) ([]*model.Withdraw, error) {
	rows, err := w.db.Query(`select order_num, sum, processed_at
								from withdraws 
								where login = $1 
								order by 2 desc;`, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdraws []*model.Withdraw
	for rows.Next() {
		var withdraw model.Withdraw
		err = rows.Scan(&withdraw.Order, &withdraw.Sum, &withdraw.ProcessedAt)
		if err != nil {
			return nil, err
		}

		withdraws = append(withdraws, &withdraw)
	}
	return withdraws, nil
}

// запрос на списание баллов
func (w *WithdrawRepository) AddWithdraw(login string, order string, sum float64, processedAt string) error {
	_, err := w.db.Exec(`INSERT INTO withdraws (order_num, sum, login, processed_at)
						 VALUES ($1, $2, $3, $4);`, order, sum, login, processedAt)
	if err != nil {
		return err
	}
	return nil
}
