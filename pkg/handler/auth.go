package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riabkovK/microgreens/internal"
)

type signUpResponse struct {
	UserId int `json:"user_id"`
}

func (h *Handler) signUp(c *fiber.Ctx) error {
	request := internal.User{}

	if err := c.BodyParser(&request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(request); err != nil {
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
		Email    string `json:"email" validate:"required,email,max=64"`
		Password string `json:"password" validate:"required,min=8,max=64"`
	}

	TokenResponse struct {
		AccessToken string `json:"access_token"`
	}
)

func (h *Handler) signIn(c *fiber.Ctx) error {
	request := signInRequest{}

	if err := c.BodyParser(&request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	token, err := h.services.Authorization.GenerateToken(request.Email, request.Password)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(TokenResponse{AccessToken: token})
}
