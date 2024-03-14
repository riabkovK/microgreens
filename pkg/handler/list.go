package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riabkovK/microgreens/internal"
	"github.com/sirupsen/logrus"
	"strconv"
)

func (h *Handler) createList(c *fiber.Ctx) error {
	id, err := c.Locals(userCtx).(int)
	if !err {
		return newErrorResponse(c, fiber.StatusInternalServerError, "user id from context is not type of int")
	}
	if id == 0 {
		return newErrorResponse(c, fiber.StatusInternalServerError, "user id not found")
	}

	request := internal.MicrogreensList{}
	if err := c.BodyParser(&request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// call service method
	logrus.Warning("i'm here. Id is - %v", id)
	return c.Status(fiber.StatusCreated).JSON(signInResponse{AccessToken: strconv.Itoa(id)})
}

func (h *Handler) getAllList(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) getListById(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) updateList(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) deleteList(c *fiber.Ctx) error {
	return nil
}
