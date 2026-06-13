package routes

import (
<<<<<<< HEAD
	"github.com/gofiber/fiber/v2"
=======
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
>>>>>>> 5497d4e04c73e5c7c543b1d5d5260635bdab344f

	"github.com/TheChoy/Cinema_Ticket_Booking/internal/handlers"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/middleware"
)

func Setup(app *fiber.App) {
<<<<<<< HEAD
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
=======

	app.Post("/login", handlers.Login)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	app.Use(middleware.CheckMiddleware)

	app.Get("/books", handlers.GetBooks)
	app.Get("/books/:id", handlers.GetBook)
	app.Post("/books", handlers.CreateBook)
}
>>>>>>> 5497d4e04c73e5c7c543b1d5d5260635bdab344f
