package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/riabkovK/microgreens/internal"
)

// Errors
var (
	errUserAlreadyExist = errors.New("the user is already exist")
	errUserNotFound     = errors.New("user not found")
	errUserNameIsNotSet = errors.New("user name is not set")
	errBadCredentials   = errors.New("email or password is incorrect")
)

type signUpResponse struct {
	UserId int `json:"user_id"`
}

func (h *Handler) signUp(c *fiber.Ctx) error {
	request := internal.User{}

	if err := c.BodyParser(&request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	id, err := h.services.Authorization.CreateUser(request)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(signUpResponse{
		UserId: id,
	})
}

type (
	signInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	signInResponse struct {
		AccessToken string `json:"access_token"`
	}
)

func (h *Handler) signIn(c *fiber.Ctx) error {
	request := signInRequest{}

	if err := c.BodyParser(&request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	token, err := h.services.Authorization.GenerateToken(request.Email, request.Password)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(signInResponse{AccessToken: token})
}
