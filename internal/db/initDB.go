package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgre driver for sql
)

// InitConnDB initiates connection to database and prepare interface
// for interaction with database
func InitConnDB(cfg Config) (*Handler, error) {
	db, err := sqlx.Open("postgres", cfg.DBAddress)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)

	err = db.Ping()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	handler := &Handler{
		DB: db,
	}

	if err := handler.prepareStatements(); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return handler, nil
}
