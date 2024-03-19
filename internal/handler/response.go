package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/riabkovK/microgreens/internal/domain"
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

type statusResponse struct {
	Status string
}

type microgreensStructures interface {
	domain.MicrogreensList | domain.MicrogreensItem | domain.MicrogreensFamily
}

type getAllResponse[T microgreensStructures] struct {
	Data []T
}
