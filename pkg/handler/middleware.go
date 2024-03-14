package handler

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *fiber.Ctx) error {
	header := c.Get(authorizationHeader)
	if header == "" {
		return newErrorResponse(c, fiber.StatusUnauthorized, "empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		return newErrorResponse(c, fiber.StatusUnauthorized, "invalid auth header")
	}

	userId, err := h.services.ParseToken(headerParts[1])
	if err != nil {
		return newErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	c.Locals(userCtx, userId)
	return c.Next()
}

func getUserId(c *fiber.Ctx) (int, error) {
	id, ok := c.Locals(userCtx).(int)
	if !ok {
		return 0, newErrorResponse(c, fiber.StatusInternalServerError, "user id from context is not type of int")
	}
	if id == 0 {
		return 0, newErrorResponse(c, fiber.StatusInternalServerError, "user id not found")
	}
	return id, nil
}
