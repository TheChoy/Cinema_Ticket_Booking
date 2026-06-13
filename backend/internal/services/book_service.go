package services

import "github.com/TheChoy/Cinema_Ticket_Booking/internal/models"

var Books []models.Book

func SeedBooks() {
	Books = append(Books,
		models.Book{
			ID:     1,
			Title:  "Spai",
			Author: "TheChoy",
		},
	)

	Books = append(Books,
		models.Book{
			ID:     2,
			Title:  "Mod",
			Author: "TheChoy",
		},
	)
}