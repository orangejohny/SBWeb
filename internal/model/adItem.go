package model

import (
	"time"
)

// TODO add tags for JSON, validator

// AdItem struct describes ad that users supposed to add
type AdItem struct {
	ID           int
	Title        string
	Price        int
	Address      string
	Owner        User
	Description  string
	CreationDate time.Time
	Tocken       string
}
