package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

type TokenManager interface {
	NewJWT(userId int, tokenTTL time.Duration) (string, error)
	Parse(accessToken string) (int, error)
	NewRefreshToken() (string, error)
}

type JWTManager struct {
	signingKey string
}

func NewJWTManager(signingKey string) (*JWTManager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &JWTManager{signingKey: signingKey}, nil
}

func (m *JWTManager) NewJWT(userId int, tokenTTL time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(tokenTTL).Unix(),
		"sub": userId,
		"iat": time.Now().Unix()})

	return token.SignedString([]byte(m.signingKey))
}

func (m *JWTManager) Parse(accessToken string) (int, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing jwt method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("error get claims from token")
	}

	// It always returns float64, so do an explicit cast to convert to int later
	userId, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("error get userId from claims")
	}

	return int(userId), err
}

func (m *JWTManager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
