package db

import "database/sql"

// Handler is used to store database connection and
// implements database interface
type Handler struct {
	DB *sql.DB

	CreateUser *sql.Stmt
	CreateAd   *sql.Stmt
	UpdateUser *sql.Stmt
	UpdateAd   *sql.Stmt
	ReadAds    *sql.Stmt
	ReadAd     *sql.Stmt
	ReadUser   *sql.Stmt
	DeleteUser *sql.Stmt
	DeleteAd   *sql.Stmt
}
