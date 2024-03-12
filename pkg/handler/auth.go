package handler

import "github.com/gofiber/fiber/v2"

func (h *Handler) signUp(c *fiber.Ctx) error {
	return c.SendString("I'm here")
}

func (h *Handler) signIn(c *fiber.Ctx) error {
	return nil
}
