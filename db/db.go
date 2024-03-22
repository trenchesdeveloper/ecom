package db

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

func NewDB(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to database")

	return db, nil
}
