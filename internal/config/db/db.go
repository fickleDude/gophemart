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

func GetDBConnection(databaseURI string) *sql.DB {
	once.Do(func() {
		dataSourceName := databaseURI
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
