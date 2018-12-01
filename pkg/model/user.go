// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"gopkg.in/guregu/null.v3/zero"
)

// User struct describes user of service
type User struct {
	ID            int64       `db:"id" json:"id" schema:"id,optional" valid:"-"`
	FirstName     string      `db:"first_name" json:"first_name" schema:"first_name,optional" valid:"utfletter,optional"` // required in DB
	LastName      string      `db:"last_name" json:"last_name" schema:"last_name,optional" valid:"utfletter,optional"`    // required in DB
	Email         string      `db:"email" json:"email" schema:"email,optional" valid:"email,optional"`                    // required in DB
	Password      string      `db:"password_hash" json:"-" schema:"password,optional" valid:"printableascii,optional"`    // required in DB
	TelNumber     zero.String `db:"telephone" json:"tel_number,omitempty" schema:"tel_number,optional" valid:"-"`         // consists of digits 1-9
	About         zero.String `db:"about" json:"about,omitempty" schema:"about,optional" valid:"-"`                       // consists of ASCII
	AvatarAddress zero.String `db:"avatar_address" json:"avatar_address,omitempty" schema:"avatar_address,optional" valid:"-"`
	RegTime       time.Time   `db:"reg_time" json:"reg_time" schema:"-" valid:"-"`
}

// TODO about should be valid UTF-8
