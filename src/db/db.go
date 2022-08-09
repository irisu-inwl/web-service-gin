package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type DBHandler struct {
	db *sql.DB
}

func configureDatabase() *sql.DB {
	db, err := sql.Open("pgx", "host=postgres port=5432 user=test dbname=testdb password=test sslmode=disable")
	if nil != err {
		panic(err)
	}
	return db
}

func NewHandler() *DBHandler {
	db := configureDatabase()
	return &DBHandler{db: db}
}

func (handler DBHandler) GetDB() *sql.DB {
	return handler.db
}
