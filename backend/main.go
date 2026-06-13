package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
<<<<<<< HEAD

	"github.com/TheChoy/Cinema_Ticket_Booking/config"
	"github.com/TheChoy/Cinema_Ticket_Booking/database"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/middleware"
=======
	"github.com/joho/godotenv"

>>>>>>> 5497d4e04c73e5c7c543b1d5d5260635bdab344f
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/routes"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/services"
)

func main() {
<<<<<<< HEAD
	config.Load()
	database.ConnectMongo()

	if err := middleware.InitFirebase(); err != nil {
		log.Fatalf("Firebase init error: %v", err)
=======

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
>>>>>>> 5497d4e04c73e5c7c543b1d5d5260635bdab344f
	}

	app := fiber.New()

	services.SeedBooks()

	routes.Setup(app)

<<<<<<< HEAD
	if err := app.Listen(":" + config.C.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
=======
	app.Listen(":8081")
>>>>>>> 5497d4e04c73e5c7c543b1d5d5260635bdab344f
}