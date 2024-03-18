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
	Create(userId int, list internal.MicrogreensList) (int, error)
	GetAll(userId int) ([]internal.MicrogreensList, error)
	GetById(userId, listId int) (internal.MicrogreensList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, request internal.UpdateMicrogreensListRequest) error
}

type MicrogreensItem interface {
	Create(listId int, microgreensItem internal.MicrogreensItem) (int, error)
	GetAll(userId, listId int) ([]internal.MicrogreensItem, error)
	GetById(userId, itemId int) (internal.MicrogreensItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, request internal.UpdateMicrogreensItemRequest) error
}

type Storage struct {
	Authorization
	MicrogreensList
	MicrogreensItem
}

func NewSQLStorage(db *sqlx.DB) *Storage {
	return &Storage{
		Authorization:   NewAuthPostgres(db),
		MicrogreensList: NewMicrogreensListPostgres(db),
		MicrogreensItem: NewMicrogreensItemPostgres(db),
	}
}
