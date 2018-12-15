// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

/*
Package model describes an internal architecture of application.

Two main data types of service are User and AdItem - they represents content of the service.
Clients operates with this two data types through the server, that's core idea of the service.

Internal architecture of API

This package contains several data types and interfaces that are used by other packages.
API uses Model struct to access database and session manager not directly but through
setted interface. It helps to make application with different technologies used in session manager and
database without changing api package.

So, both db and sessionmanager packages must contain structure that implements interfaces from this package.

Data types

User - a human that can sign up, sign in, logout, create ads and etc.
Characteristics:
	id	                  unique identificator of user in database
	first name            first name of user
	last name	            last name of user
	email                 email of user that must be unique
	password              password of user that stored in database in hashed state
	telephone number      telephone number of user
	about                 some additional information about user

Ad - main content of application
Characteristics:
	id                    unique identificator of ad in database
	title                 title of ad
	price                 price of ad
	country               country where ad is situated
	city                  city where ad is situated
	subway station        station where ad is situated
	images                images of ad
	owner                 user which created this ad
	description           additional information about this ad
	creation time         time when this ad was created
*/
package model

// Model is a struct that contains DB, IM and SM interfaces. Such project model allows
// to use different database and session manager implementation without changing business-logic.
// Model is used by API handlers.
type Model struct {
	DB
	SM
	IM
}

// New creates Model structure from object that implements DB and SM interfaces.
func New(db DB, sm SM, im IM) *Model {
	return &Model{
		DB: db,
		SM: sm,
		IM: im,
	}
}
