package service

import (
	"github.com/riabkovK/microgreens/internal/domain"
	"github.com/riabkovK/microgreens/pkg/auth"
	"github.com/riabkovK/microgreens/pkg/hash"
	"time"

	"github.com/riabkovK/microgreens/internal/storage"
)

type AuthService struct {
	storage    storage.Authorization
	hasher     hash.PasswordHasher
	jwtManager auth.TokenManager

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthService(storages storage.Authorization, hasher hash.PasswordHasher, jwtManager auth.TokenManager,
	accessTokenTTL, refreshTokenTTL time.Duration) *AuthService {
	return &AuthService{
		storage:         storages,
		hasher:          hasher,
		jwtManager:      jwtManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL}
}

type UserSignUpRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserSignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (as *AuthService) SignUp(request UserSignUpRequest) (int, error) {
	hashedPassword := as.hasher.Hash(request.Password)
	return as.storage.CreateUser(domain.User{
		Email:    request.Email,
		Name:     request.Name,
		Username: request.Username,
		Password: hashedPassword,
	})
}

func (as *AuthService) SingIn(request UserSignInRequest) (domain.TokensResponse, error) {
	passwordHash := as.hasher.Hash(request.Password)
	user, err := as.storage.GetByCredentials(request.Email, passwordHash)
	if err != nil {
		return domain.TokensResponse{}, nil
	}

	return as.createSession(user.Id)
}

func (as *AuthService) RefreshTokens(refreshToken string) (domain.TokensResponse, error) {
	user, err := as.storage.GetByRefreshToken(refreshToken)
	if err != nil {
		return domain.TokensResponse{}, err
	}

	return as.createSession(user.Id)
}

func (as *AuthService) createSession(userId int) (domain.TokensResponse, error) {
	var (
		res domain.TokensResponse
		err error
	)

	res.AccessToken, err = as.jwtManager.NewJWT(userId, as.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = as.jwtManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(as.refreshTokenTTL),
	}

	err = as.storage.SetSession(userId, session)

	return res, err
}
