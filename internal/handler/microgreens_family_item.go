package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/riabkovK/microgreens/internal/domain"
)

func (h *Handler) createFamily(c *fiber.Ctx) error {
	request := domain.MicrogreensFamilyRequest{}
	if err := c.BodyParser(&request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	id, err := h.services.MicrogreensFamily.Create(request)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(responseWithId{Id: id})
}

func (h *Handler) getAllFamilies(c *fiber.Ctx) error {
	families, err := h.services.MicrogreensFamily.GetAll()
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(getAllResponse[domain.MicrogreensFamily]{Data: families})
}

func (h *Handler) getFamilyById(c *fiber.Ctx) error {
	familyId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, "invalid id param")
	}

	family, err := h.services.MicrogreensFamily.GetById(familyId)
	if err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(family)
}

func (h *Handler) updateFamily(c *fiber.Ctx) error {
	familyId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, "invalid id param")
	}

	request := domain.UpdateMicrogreensFamilyRequest{}
	if err := c.BodyParser(&request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(request); err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.services.MicrogreensFamily.Update(familyId, request); err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(statusResponse{Status: "ok"})
}

func (h *Handler) deleteFamily(c *fiber.Ctx) error {
	familyId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return newErrorResponse(c, fiber.StatusBadRequest, "invalid id param")
	}

	rows, err := h.services.MicrogreensFamily.Delete(familyId)
	if err != nil {
		return newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if rows == 0 {
		return c.SendStatus(fiber.StatusNoContent)
	}

	return c.Status(fiber.StatusOK).JSON(statusResponse{Status: "Successfully removed"})
}
