package service

import "github.com/riabkovK/microgreens/internal/storage"

type Authorization interface {
}

type MicrogreensList interface {
}

type MicrogreensItem interface {
}

type Service struct {
	Authorization
	MicrogreensList
	MicrogreensItem
}

func NewService(storages *storage.Storage) *Service {
	return &Service{}
}
