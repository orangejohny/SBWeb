package model

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

// User struct describes user of service
type User struct {
	ID        int64       `db:"id" json:"id" schema:"id,optional" valid:"-"`
	FirstName string      `db:"first_name" json:"first_name" schema:"first_name,optional" valid:"utfletter,optional"`
	LastName  string      `db:"last_name" json:"last_name" schema:"last_name,optional" valid:"utfletter,optional"`
	Email     string      `db:"email" json:"email" schema:"email,optional" valid:"email,required" `
	Password  string      `db:"password_hash" json:"-" schema:"password,optional" valid:",required"`
	TelNumber null.String `db:"telephone" json:"tel_number,omitempty" schema:"tel_number,optional" valid:"-"`
	About     null.String `db:"about" json:"about,omitempty" schema:"about,optional" valid:",optional"`
	RegTime   time.Time   `db:"reg_time" json:"reg_time" schema:"-" valid:"-"`
}
