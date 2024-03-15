package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riabkovK/microgreens/internal"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	logrus.Error(message)
	return c.Status(statusCode).JSON(errorResponse{Message: message})
}

type responseWithId struct {
	Id int `json:"id"`
}

type getAllListResponse struct {
	Data []internal.MicrogreensList `json:"data"`
}

type getAllItemsResponse struct {
	Data []internal.MicrogreensItem `json:"data"`
}

type statusResponse struct {
	Status string
}
