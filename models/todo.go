package models

import "time"

type Todo struct {
	ID         int `pg:",pk,unique,notnull"`
	Name       string
	IsComplete bool
	IsDeleted  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time

	CreatedBy int
	UpdatedBy int
}

type TodoInput struct {
	Name      string
	CreatedBy int
}
