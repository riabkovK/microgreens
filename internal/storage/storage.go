package storage

import (
	"github.com/jmoiron/sqlx"

	"github.com/riabkovK/microgreens/internal/domain"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GetByCredentials(email, passwordHash string) (domain.User, error)
	GetByRefreshToken(refreshToken string) (domain.User, error)
	SetSession(userId int, session domain.Session) error
}

type MicrogreensList interface {
	Create(userId int, list domain.MicrogreensListRequest) (int, error)
	GetAll(userId int) ([]domain.MicrogreensList, error)
	GetById(userId, listId int) (domain.MicrogreensList, error)
	Delete(userId, listId int) (int, error)
	Update(userId, listId int, request domain.UpdateMicrogreensListRequest) error
}

type MicrogreensItem interface {
	Create(listId int, microgreensItem domain.MicrogreensItemRequest) (int, error)
	GetAll(userId, listId int) ([]domain.MicrogreensItem, error)
	GetById(userId, itemId int) (domain.MicrogreensItem, error)
	Delete(userId, itemId int) (int, error)
	Update(userId, itemId int, request domain.UpdateMicrogreensItemRequest) error
}

type MicrogreensFamily interface {
	Create(family domain.MicrogreensFamilyRequest) (int, error)
	GetAll() ([]domain.MicrogreensFamily, error)
	GetById(familyId int) (domain.MicrogreensFamily, error)
	Delete(familyId int) (int, error)
	Update(familyId int, request domain.UpdateMicrogreensFamilyRequest) error
}

type Storage struct {
	Authorization
	MicrogreensList
	MicrogreensItem
	MicrogreensFamily
}

func NewSQLStorage(db *sqlx.DB) *Storage {
	return &Storage{
		Authorization:     NewAuthPostgres(db),
		MicrogreensList:   NewMicrogreensListPostgres(db),
		MicrogreensItem:   NewMicrogreensItemPostgres(db),
		MicrogreensFamily: NewMicrogreensFamilyPostgres(db),
	}
}
