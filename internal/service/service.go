package service

import (
	"time"

	"github.com/riabkovK/microgreens/internal/domain"
	"github.com/riabkovK/microgreens/internal/storage"
	"github.com/riabkovK/microgreens/pkg/auth"
	"github.com/riabkovK/microgreens/pkg/hash"
)

type Authorization interface {
	SignUp(request UserSignUpRequest) (int, error)
	SingIn(request UserSignInRequest) (domain.TokensResponse, error)
	RefreshTokens(refreshToken string) (domain.TokensResponse, error)
}

type MicrogreensList interface {
	Create(userId int, list domain.MicrogreensListRequest) (int, error)
	GetAll(userId int) ([]domain.MicrogreensList, error)
	GetById(userId, listId int) (domain.MicrogreensList, error)
	Delete(userId, listId int) (int, error)
	Update(userId, listId int, request domain.UpdateMicrogreensListRequest) error
}

type MicrogreensItem interface {
	Create(userId, listId int, microgreensItem domain.MicrogreensItemRequest) (int, error)
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

type Service struct {
	Authorization
	MicrogreensList
	MicrogreensItem
	MicrogreensFamily
}

// Dependencies

type Deps struct {
	Storages        *storage.Storage
	Hasher          hash.PasswordHasher
	JWTManager      auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewService(deps Deps) *Service {
	return &Service{
		Authorization:     NewAuthService(deps.Storages.Authorization, deps.Hasher, deps.JWTManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
		MicrogreensList:   NewMicrogreensListService(deps.Storages.MicrogreensList),
		MicrogreensItem:   NewMicrogreensItemService(deps.Storages.MicrogreensItem, deps.Storages.MicrogreensList),
		MicrogreensFamily: NewMicrogreensFamilyService(deps.Storages.MicrogreensFamily),
	}
}
