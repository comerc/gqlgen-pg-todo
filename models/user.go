package models

import "time"

type User struct {
	ID    int    `pg:",pk,unique,notnull"`
	Email string `pg:",unique,notnull"`

	FirstName string
	LastName  string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserInput struct {
	Email     string
	FirstName string
	LastName  string
}
