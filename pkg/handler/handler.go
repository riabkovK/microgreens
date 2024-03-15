package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/riabkovK/microgreens/pkg/service"
)

type Handler struct {
	services *service.Service
	validate *validator.Validate
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services,
		validate: validator.New(validator.WithRequiredStructEnabled())}
}

func (h *Handler) SetupRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	{
		auth.Post("sign-up", h.signUp)
		auth.Post("sign-in", h.signIn)
	}

	api := app.Group("/api", h.userIdentity)

	{
		lists := api.Group("lists")
		{
			lists.Post("/", h.createList)
			lists.Get("/", h.getAllLists)
			lists.Get("/:id", h.getListById)
			lists.Put("/:id", h.updateList)
			lists.Delete("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.Post("/", h.createItem)
				items.Get("/", h.getAllItems)
			}
		}

		items := api.Group("items")
		{
			items.Get("/:id", h.getItemById)
			items.Put("/:id", h.updateItem)
			items.Delete("/:id", h.deleteItem)
		}

		families := api.Group("families")
		{
			families.Post("/")
			families.Get("/")
			families.Get("/:id")
			families.Put("/:id")
			families.Delete("/:id")
		}

	}
}
