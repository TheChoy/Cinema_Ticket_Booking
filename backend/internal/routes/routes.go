package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"

	"github.com/TheChoy/Cinema_Ticket_Booking/internal/handlers"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/middleware"
)

func Setup(app *fiber.App) {

	app.Post("/login", handlers.Login)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	app.Use(middleware.CheckMiddleware)

	app.Get("/books", handlers.GetBooks)
	app.Get("/books/:id", handlers.GetBook)
	app.Post("/books", handlers.CreateBook)
}