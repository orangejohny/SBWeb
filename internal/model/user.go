package model

import "time"

// TODO add tags

// User struct describes user of service
type User struct {
	ID        int
	FirstName string
	LastName  string
	Nickname  string
	Email     string
	TelNumber string
	About     string
	RegDate   time.Time
}
