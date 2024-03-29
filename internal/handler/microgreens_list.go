package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/riabkovK/microgreens/internal/domain"
)

func (h *Handler) createList(c *fiber.Ctx) error {
	userId, err := getUserId(c)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	request := domain.MicrogreensListRequest{}
	if err := c.BodyParser(&request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	id, err := h.services.MicrogreensList.Create(userId, request)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(responseWithId{Id: id})
}

func (h *Handler) getAllLists(c *fiber.Ctx) error {
	userId, err := getUserId(c)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	lists, err := h.services.MicrogreensList.GetAll(userId)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(getAllResponse[domain.MicrogreensList]{Data: lists})
}

func (h *Handler) getListById(c *fiber.Ctx) error {
	userId, err := getUserId(c)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	microgreensListId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, "invalid id param")
	}

	list, err := h.services.MicrogreensList.GetById(userId, microgreensListId)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(list)
}

func (h *Handler) updateList(c *fiber.Ctx) error {
	userId, err := getUserId(c)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	microgreensListId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, "invalid id param")
	}

	request := domain.UpdateMicrogreensListRequest{}
	if err := c.BodyParser(&request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.services.MicrogreensList.Update(userId, microgreensListId, request); err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(statusResponse{Status: "ok"})
}

func (h *Handler) deleteList(c *fiber.Ctx) error {
	userId, err := getUserId(c)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	microgreensListId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, "invalid id param")
	}

	rows, err := h.services.MicrogreensList.Delete(userId, microgreensListId)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if rows == 0 {
		return c.SendStatus(fiber.StatusNoContent)
	}

	return c.Status(fiber.StatusOK).JSON(statusResponse{Status: "Successfully removed"})
}
