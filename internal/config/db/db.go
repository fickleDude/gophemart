package db

import (
	"database/sql"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	instance *sql.DB
	once     sync.Once
)

func GetDBConnection() *sql.DB {
	once.Do(func() {
		dataSourceName := "host=localhost port=5433 user=postgres password=postgres dbname=gophermart sslmode=disable"
		db, err := sql.Open("pgx", dataSourceName)
		if err != nil {
			panic(err.Error())
		}
		//init database
		instance = db
	})
	return instance
}

func CloseDBConnection() error {
	return instance.Close()
}
