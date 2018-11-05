package model

import (
	"time"
)

// User struct describes user of service
type User struct {
	ID        int64     `db:"id" json:"id" schema:"id,optional" valid:"-"`
	FirstName string    `db:"first_name" json:"first_name" schema:"first_name,optional" valid:"alpha,optional"`
	LastName  string    `db:"last_name" json:"last_name" schema:"last_name,optional" valid:"alpha,optional"`
	Email     string    `db:"email" json:"email" schema:"email,required" valid:"email,required" `
	Password  string    `db:"password_hash" json:"-" schema:"password,required" valid:",required"`
	TelNumber string    `db:"telephone" json:"tel_number,omitempty" schema:"tel_number,optional" valid:"-"`
	About     string    `db:"about" json:"about,omitempty" schema:"about,optional" valid:"utfletternum,optional"`
	RegTime   time.Time `db:"reg_time" json:"reg_time" schema:"-" valid:"-"`
}
