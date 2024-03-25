package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/riabkovK/microgreens/internal/service"
	"github.com/riabkovK/microgreens/pkg/auth"
)

type Handler struct {
	services     *service.Service
	validate     *validator.Validate
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Service, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		validate:     validator.New(validator.WithRequiredStructEnabled()),
		tokenManager: tokenManager}
}

func (h *Handler) SetupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Hello, that's microgreens")
	})

	auth := app.Group("/auth")
	{
		auth.Post("sign-up", h.userSignUp)
		auth.Post("sign-in", h.userSignIn)
		auth.Post("refresh", h.userRefresh)
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
			families.Post("/", h.createFamily)
			families.Get("/", h.getAllFamilies)
			families.Get("/:id", h.getFamilyById)
			families.Put("/:id", h.updateFamily)
			families.Delete("/:id", h.deleteFamily)
		}

	}
}
