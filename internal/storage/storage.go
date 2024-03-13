package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/riabkovK/microgreens/internal"
)

type Authorization interface {
	CreateUser(user internal.User) (int, error)
	GetUser(email, password string) (internal.User, error)
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
	return &Storage{
		Authorization: NewAuthPostgres(db),
	}
}
