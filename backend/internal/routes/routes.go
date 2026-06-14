package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"


	"github.com/TheChoy/Cinema_Ticket_Booking/internal/handlers"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/middleware"
)

func Setup(app *fiber.App) {
	// Public
	app.Post("/auth/login", middleware.AuthMiddleware, handlers.Register)
	app.Get("/movies", handlers.GetMovies)
	app.Get("/movies/:id", handlers.GetMovieByID)
	app.Get("/showtimes", handlers.GetShowtimes)
	app.Get("/seats", handlers.GetSeats)

	// WebSocket - real-time seat status
	app.Get("/ws/showtimes/:showtime_id", websocket.New(handlers.WSSeatStatus))

	// Protected
	protected := app.Group("/", middleware.AuthMiddleware)
	protected.Get("/users/me", handlers.GetMe)
	protected.Put("/users/me", handlers.UpdateMe)
	protected.Post("/bookings", handlers.CreateBooking)
	protected.Get("/bookings/me", handlers.GetMyBookings)
	protected.Get("/bookings/:id", handlers.GetBookingByID)
	protected.Put("/bookings/:id/pay", handlers.PayBooking)



	// Admin
	admin := app.Group("/admin", middleware.AuthMiddleware, middleware.AdminOnly)
	admin.Post("/movies", handlers.CreateMovie)
	admin.Put("/movies/:id", handlers.UpdateMovie)
	admin.Delete("/movies/:id", handlers.DeleteMovie)
	admin.Post("/showtimes", handlers.CreateShowtime)
	admin.Put("/showtimes/:id", handlers.UpdateShowtime)
	admin.Delete("/showtimes/:id", handlers.DeleteShowtime)
	admin.Post("/seats/generate", handlers.GenerateSeats)
	admin.Get("/users", handlers.GetUsers)
	admin.Put("/users/:id/role", handlers.UpdateUserRole)
	admin.Get("/bookings", handlers.GetAllBookings)
	admin.Put("/bookings/:id", handlers.UpdateBooking)
	admin.Delete("/bookings/:id", handlers.CancelBooking)
}