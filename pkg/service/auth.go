package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/riabkovK/microgreens/internal"
	"github.com/riabkovK/microgreens/internal/storage"

	"github.com/golang-jwt/jwt/v5"
)

const (
	passwordSalt  = "t3R/i)96DGg{a{d2"
	jwtSigningKey = ("4>p4UvtV>}46#8hwu%1lF")
	tokenTTL      = 12 * time.Hour
)

type tokenClaims struct {
	jwt.MapClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	storage storage.Authorization
}

func NewAuthService(storages storage.Authorization) *AuthService {
	return &AuthService{storage: storages}
}

func (as *AuthService) CreateUser(user internal.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return as.storage.CreateUser(user)
}

func (as *AuthService) GenerateToken(email, password string) (string, error) {
	user, err := as.storage.GetUser(email, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.MapClaims{
			"exp": time.Now().Add(tokenTTL).Unix(),
			"iat": time.Now().Unix()},
		user.Id,
	})

	return token.SignedString([]byte(jwtSigningKey))
}

func (as *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing jwt method")
		}
		return []byte(jwtSigningKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hashedPassword := sha256.Sum256([]byte(password + passwordSalt))

	return fmt.Sprintf("%x", hashedPassword)
}
