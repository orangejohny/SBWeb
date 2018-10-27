package model

import (
	"time"
)

// TODO add tags

// User struct describes user of service
type User struct {
	ID        int       `db:"id" json:"id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password_hash" json:"-"`
	TelNumber string    `db:"telephone" json:"tel_number,omitempty"`
	About     string    `db:"about" json:"about,omitempty"`
	RegTime   time.Time `db:"reg_time" json:"reg_time"`
}
