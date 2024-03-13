package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Error struct {
	Message string `json:"message"`
}

func newErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	logrus.Error(message)
	return c.Status(statusCode).JSON(Error{Message: message})
}
