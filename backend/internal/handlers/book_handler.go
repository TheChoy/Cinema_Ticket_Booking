package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/TheChoy/Cinema_Ticket_Booking/internal/models"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/services"
)

func GetBooks(c *fiber.Ctx) error {
	return c.JSON(services.Books)
}

func GetBook(c *fiber.Ctx) error {

	bookID, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for _, book := range services.Books {
		if book.ID == bookID {
			return c.JSON(book)
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func CreateBook(c *fiber.Ctx) error {

	book := new(models.Book)

	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	services.Books = append(services.Books, *book)

	return c.JSON(book)
}