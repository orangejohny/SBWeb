package db

import (
	"github.com/jmoiron/sqlx"
)

// Handler is used to store database connection and
// implements database interface needed by API
type Handler struct {
	DB *sqlx.DB

	CreateUser        *sqlx.NamedStmt
	CreateAd          *sqlx.NamedStmt
	UpdateUser        *sqlx.NamedStmt
	UpdateAd          *sqlx.NamedStmt
	ReadAds           *sqlx.NamedStmt
	SearchAds         *sqlx.NamedStmt
	ReadAdsOfUser     *sqlx.Stmt
	ReadAd            *sqlx.Stmt
	ReadUserWithID    *sqlx.Stmt
	ReadUserWithEmail *sqlx.Stmt
	DeleteUser        *sqlx.Stmt
	DeleteAd          *sqlx.Stmt
}
