package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/riabkovK/microgreens/internal/domain"
	"github.com/riabkovK/microgreens/internal/service"
)

type userSignUpResponse struct {
	UserId int `json:"user_id"`
}

type signUpRequest struct {
	Email    string `json:"email" validate:"required,email,max=64"`
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required,max=25"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func (h *Handler) userSignUp(c *fiber.Ctx) error {
	request := signUpRequest{}

	if err := c.BodyParser(&request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	id, err := h.services.Authorization.SignUp(service.UserSignUpRequest{
		Email:    request.Email,
		Name:     request.Name,
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(userSignUpResponse{
		UserId: id,
	})
}

type userSignInRequest struct {
	Email    string `json:"email" validate:"required,email,max=64"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func (h *Handler) userSignIn(c *fiber.Ctx) error {
	request := userSignInRequest{}

	if err := c.BodyParser(&request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	token, err := h.services.Authorization.SingIn(service.UserSignInRequest{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(domain.TokensResponse{
		RefreshToken: token.RefreshToken,
		AccessToken:  token.AccessToken,
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *Handler) userRefresh(c *fiber.Ctx) error {
	request := refreshRequest{}

	if err := c.BodyParser(&request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	res, err := h.services.Authorization.RefreshTokens(request.RefreshToken)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(domain.TokensResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}
