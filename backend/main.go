package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/TheChoy/Cinema_Ticket_Booking/internal/routes"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/services"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	services.SeedBooks()

	routes.Setup(app)

	app.Listen(":8081")
}