package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *fiber.Ctx) error {
	userId, err := h.parseAuthHeader(c)
	if err != nil {
		return newErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	c.Locals(userCtx, userId)
	return c.Next()
}

func (h *Handler) parseAuthHeader(c *fiber.Ctx) (int, error) {
	header := c.Get(authorizationHeader)
	if header == "" {
		return 0, errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return 0, errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return 0, errors.New("token is empty")
	}

	return h.tokenManager.Parse(headerParts[1])
}

func getUserId(c *fiber.Ctx) (int, error) {
	return getIdByContext(c, userCtx)
}

func getIdByContext(c *fiber.Ctx, context string) (int, error) {
	id, ok := c.Locals(context).(int)
	if !ok {
		return 0, errors.New("id from context is invalid type")
	}
	if id == 0 {
		return 0, errors.New("id not found")
	}
	return id, nil
}
