package storage

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/riabkovK/microgreens/internal"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// Context key for JWT
const (
	contextKeyUser = "user"
)

// Secret keys
var (
	jwtSecretKey = []byte("4>p4UvtV>}46")
	passwordSalt = "t3R/i)96DGg{"
)

// Errors
var (
	errUserAlreadyExist = errors.New("the user is already exist")
	errUserNotFound     = errors.New("user not found")
	errUserNameIsNotSet = errors.New("user name is not set")
	errBadCredentials   = errors.New("email or password is incorrect")
)

type (
	SignUpRequest struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	SignInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	SignInResponse struct {
		AccessToken string `json:"access_token"`
	}

	ProfileResponse struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
)

type AuthStorage struct {
	mu    sync.Mutex
	users map[string]internal.User
}

// AuthStorage should be only in one instance
var (
	authStorageOnce     sync.Once
	authStorageInstance *AuthStorage = nil
)

func (as *AuthStorage) addUser(email, name, password string) error {
	if email == "" {
		return errBadCredentials
	}

	if name == "" {
		return errUserNameIsNotSet
	}

	if password == "" {
		return errBadCredentials
	}

	as.mu.Lock()
	as.users[email] = internal.User{
		Email:        email,
		Name:         name,
		PasswordHash: password}
	as.mu.Unlock()

	return nil
}

func (as *AuthStorage) getUser(email string) (*internal.User, error) {
	as.mu.Lock()
	defer as.mu.Unlock()
	user, exists := as.users[email]
	if !exists {
		return nil, errUserNotFound
	}

	return &user, nil
}

func (as *AuthStorage) updateUser(email, name string) (*internal.User, error) {
	return nil, nil
}

func (as *AuthStorage) deleteUser(email string) error {
	return nil
}

func NewAuthStorage() *AuthStorage {
	authStorageOnce.Do(func() {
		authStorageInstance = &AuthStorage{users: make(map[string]internal.User)}
	})
	return authStorageInstance
}

type AuthenticationHandler interface {
	SignUp(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
}

type AuthenticatedUserHandler interface {
	GetProfile(c *fiber.Ctx) error
}

type AuthHandler struct {
	storage *AuthStorage
}

func NewAuthHandler(storage *AuthStorage) AuthenticationHandler {
	return &AuthHandler{storage: storage}
}

func (ah *AuthHandler) SignUp(c *fiber.Ctx) error {
	regReq := SignUpRequest{}
	if err := c.BodyParser(&regReq); err != nil {
		logrus.WithError(err).Error("SignUp body parser")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, _ := ah.storage.getUser(regReq.Email)
	if user != nil {
		logrus.WithError(errUserAlreadyExist).Error("SignUp user conflict")
		return c.SendStatus(fiber.StatusConflict)
	}

	hashedPassword := hashPassword(regReq.Password)
	err := ah.storage.addUser(regReq.Email, regReq.Name, hashedPassword)
	if err != nil {
		logrus.WithError(err).Error("SignUp add user to storage")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (ah *AuthHandler) SignIn(c *fiber.Ctx) error {
	regReq := SignInRequest{}
	if err := c.BodyParser(&regReq); err != nil {
		logrus.WithError(err).Error("SignIn body parser")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	ah.storage.mu.Lock()
	user, exists := ah.storage.users[regReq.Email]
	ah.storage.mu.Unlock()

	if !exists {
		logrus.Error(errBadCredentials)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	hashedPassword := hashPassword(regReq.Password)
	if user.PasswordHash != hashedPassword {
		logrus.Error(errBadCredentials)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	payload := jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		logrus.WithError(err).Error("JWT token signing")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(SignInResponse{AccessToken: t})
}

type UserHandler struct {
	storage *AuthStorage
}

func NewAuthenticatedUserHandler(storage *AuthStorage) AuthenticatedUserHandler {
	return &UserHandler{storage: storage}
}

func jwtPayloadFromRequest(c *fiber.Ctx) (jwt.MapClaims, bool) {
	jwtToken, ok := c.Context().Value(contextKeyUser).(*jwt.Token)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"jwt_token_context_value": c.Context().Value(contextKeyUser),
		}).Error("wrong type of JWT token in context")
	}

	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"jwt_token_claims": jwtToken.Claims,
		}).Error("wrong type of JWT token claims")

		return nil, false
	}

	return payload, true
}

func (uh *UserHandler) GetProfile(c *fiber.Ctx) error {
	jwtPayload, ok := jwtPayloadFromRequest(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	uh.storage.mu.Lock()
	userInfo, ok := uh.storage.users[jwtPayload["sub"].(string)]
	if !ok {
		logrus.Error(errUserNotFound)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	uh.storage.mu.Unlock()

	return c.JSON(ProfileResponse{
		Email: userInfo.Email,
		Name:  userInfo.Name,
	})
}

func hashPassword(password string) string {
	hashedPassword := sha256.Sum256([]byte(password + passwordSalt))
	return fmt.Sprintf("%x", hashedPassword)
}
