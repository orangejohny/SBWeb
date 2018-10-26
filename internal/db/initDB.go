package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // postgre driver for sql
)

// InitConnDB initiates connection to database and prepare interface
// for interaction with database
func InitConnDB(cfg Config) (*Handler, error) {
	db, err := sql.Open("postgres", cfg.DBAddress)
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

// prepareStateents just preapares SQL code
func (h *Handler) prepareStatements() (err error) {
	if h.ReadAds, err = h.DB.Prepare(
		"SELECT * FROM ads WHERE id > $1 AND id < $1 + $2",
	); err != nil {
		return err
	}

	if h.ReadAd, err = h.DB.Prepare(
		"SELECT * FROM ads WHERE id=$1",
	); err != nil {
		return err
	}

	if h.ReadUser, err = h.DB.Prepare(
		"SELECT * FROM users WHERE id=$1",
	); err != nil {
		return err
	}

	if h.CreateUser, err = h.DB.Prepare(
		`INSERT INTO users
			(first_name, last_name, email)
			VALUES
			($1, $2, $3)`,
	); err != nil {
		return err
	}

	if h.CreateAd, err = h.DB.Prepare(
		`INSERT INTO ads
			(title, owner_ad, description)
			VALUES
			($1, $2)`,
	); err != nil {
		return err
	}

	if h.UpdateUser, err = h.DB.Prepare(
		`UPDATE users SET
			first_name=$1,
			last_name=$2,
			email=$3
			WHERE id=$4`,
	); err != nil {
		return err
	}

	if h.UpdateAd, err = h.DB.Prepare(
		`UPDATE ads SET
			title=$1,
			description=$2
			WHERE id=$3`,
	); err != nil {
		return err
	}

	if h.DeleteUser, err = h.DB.Prepare(
		`DELETE FROM users WHERE id=$1`,
	); err != nil {
		return err
	}

	if h.DeleteAd, err = h.DB.Prepare(
		`DELETE FROM ads WHERE id=$2`,
	); err != nil {
		return err
	}

	return nil
}
