package database

import (
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 5

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
		Plan: Plan{},
	}
}

type Models struct {
	User User
	Plan Plan
}

type InitPlan struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}
