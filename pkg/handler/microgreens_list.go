package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riabkovK/microgreens/internal"
)

type response struct {
	Id int `json:"id"`
}

func (h *Handler) createList(c *fiber.Ctx) error {
	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	request := internal.MicrogreensList{}
	if err := c.BodyParser(&request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	id, err := h.services.MicrogreensList.Create(userId, request)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response{Id: id})
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