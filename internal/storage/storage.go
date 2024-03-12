package storage

import "github.com/jmoiron/sqlx"

type Authorization interface {
}

type MicrogreensList interface {
}

type MicrogreensItem interface {
}

type Storage struct {
	Authorization
	MicrogreensList
	MicrogreensItem
}

func NewSQLStorage(db *sqlx.DB) *Storage {
	return &Storage{}
}
