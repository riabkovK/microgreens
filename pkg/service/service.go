package service

import (
	"github.com/riabkovK/microgreens/internal"
	"github.com/riabkovK/microgreens/internal/storage"
)

type Authorization interface {
	CreateUser(user internal.User) (int, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type MicrogreensList interface {
	Create(userId int, list internal.MicrogreensList) (int, error)
	GetAll(userId int) ([]internal.MicrogreensList, error)
	GetById(userId, listId int) (internal.MicrogreensList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, request internal.UpdateMicrogreensListRequest) error
}

type MicrogreensItem interface {
}

type Service struct {
	Authorization
	MicrogreensList
	MicrogreensItem
}

func NewService(storages *storage.Storage) *Service {
	return &Service{
		Authorization:   NewAuthService(storages.Authorization),
		MicrogreensList: NewMicrogreensListService(storages.MicrogreensList),
	}
}
