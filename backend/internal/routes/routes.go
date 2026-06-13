package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/TheChoy/Cinema_Ticket_Booking/internal/handlers"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/middleware"
)

func Setup(app *fiber.App) {
	// Public
	app.Post("/auth/login", middleware.AuthMiddleware, handlers.Register)

	// Protected
	protected := app.Group("/", middleware.AuthMiddleware)
	protected.Get("/books", handlers.GetBooks)
	protected.Get("/books/:id", handlers.GetBook)

	// Admin
	admin := app.Group("/admin", middleware.AuthMiddleware, middleware.AdminOnly)
	admin.Post("/books", handlers.CreateBook)
}
