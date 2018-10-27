package db

import (
	"github.com/jmoiron/sqlx"
)

// Handler is used to store database connection and
// implements database interface
type Handler struct {
	DB *sqlx.DB

	CreateUser        *sqlx.NamedStmt
	CreateAd          *sqlx.NamedStmt
	UpdateUser        *sqlx.NamedStmt
	UpdateAd          *sqlx.NamedStmt
	ReadAds           *sqlx.Stmt
	ReadAd            *sqlx.Stmt
	ReadUserWithID    *sqlx.Stmt
	ReadUserWithEmail *sqlx.Stmt
	DeleteUser        *sqlx.Stmt
	DeleteAd          *sqlx.Stmt
}
