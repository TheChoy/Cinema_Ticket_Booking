package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/TheChoy/Cinema_Ticket_Booking/config"
	"github.com/TheChoy/Cinema_Ticket_Booking/database"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/middleware"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/routes"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/services"
)

func main() {
	config.Load()
	database.ConnectMongo()

	if err := middleware.InitFirebase(); err != nil {
		log.Fatalf("Firebase init error: %v", err)
	}

	app := fiber.New()

	services.SeedBooks()

	routes.Setup(app)

	if err := app.Listen(":" + config.C.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}