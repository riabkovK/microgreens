package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/riabkovK/microgreens/internal/domain"
)

func (h *Handler) createItem(c *fiber.Ctx) error {
	userId, err := getUserId(c)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	microgreensListId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, "invalid id param")
	}

	var request domain.MicrogreensItemRequest
	if err := c.BodyParser(&request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	id, err := h.services.MicrogreensItem.Create(userId, microgreensListId, request)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(responseWithId{Id: id})
}

func (h *Handler) getAllItems(c *fiber.Ctx) error {
	userId, err := getUserId(c)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	microgreensListId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, "invalid id param")
	}

	items, err := h.services.MicrogreensItem.GetAll(userId, microgreensListId)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(getAllResponse[domain.MicrogreensItem]{Data: items})
}

func (h *Handler) getItemById(c *fiber.Ctx) error {
	userId, err := getUserId(c)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	itemId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, "invalid id param")
	}

	item, err := h.services.MicrogreensItem.GetById(userId, itemId)
	if err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(item)
}

func (h *Handler) updateItem(c *fiber.Ctx) error {
	userId, err := getUserId(c)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	microgreensListId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, "invalid id param")
	}

	request := domain.UpdateMicrogreensItemRequest{}
	if err := c.BodyParser(&request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.services.MicrogreensItem.Update(userId, microgreensListId, request); err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(statusResponse{Status: "ok"})
}

func (h *Handler) deleteItem(c *fiber.Ctx) error {
	userId, err := getUserId(c)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	itemId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, "invalid id param")
	}

	rows, err := h.services.MicrogreensItem.Delete(userId, itemId)
	if err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if rows == 0 {
		return c.SendStatus(fiber.StatusNoContent)
	}

	return c.Status(fiber.StatusOK).JSON(statusResponse{Status: "Successfully removed"})
}
