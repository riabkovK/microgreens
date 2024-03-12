package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riabkovK/microgreens/pkg/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{service: services}
}

func (h *Handler) SetupRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	{
		auth.Post("sign-up", h.signUp)
		auth.Post("sign-in", h.signIn)
	}

	api := app.Group("/api")

	{
		lists := api.Group("lists")
		{
			lists.Post("/", h.createList)
			lists.Get("/", h.getAllList)
			lists.Get("/:id", h.getListById)
			lists.Put("/:id", h.updateList)
			lists.Delete("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.Post("/", h.createItem)
				items.Get("/", h.getAllItem)
				items.Get("/:item_id", h.getItemById)
				items.Put("/:item_id", h.updateItem)
				items.Delete("/:item_id", h.deleteItem)
			}
		}
	}
}
