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
	return nil
}
