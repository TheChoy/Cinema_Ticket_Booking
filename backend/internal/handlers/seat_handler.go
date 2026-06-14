package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/TheChoy/Cinema_Ticket_Booking/database"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/models"
)

func GetSeats(c *fiber.Ctx) error {
	col := database.DB.Collection("seats")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	showtimeID := c.Query("showtime_id")
	if showtimeID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	id, err := primitive.ObjectIDFromHex(showtimeID)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	cursor, err := col.Find(ctx, bson.M{"showtime_id": id})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer cursor.Close(ctx)

	var seats []models.Seat
	if err := cursor.All(ctx, &seats); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(seats)
}

func GenerateSeats(c *fiber.Ctx) error {
	body := struct {
		ShowtimeID string   `json:"showtime_id"`
		Rows       []string `json:"rows"`    // เช่น ["A","B","C"]
		SeatsPerRow int     `json:"seats_per_row"` // เช่น 10
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	showtimeID, err := primitive.ObjectIDFromHex(body.ShowtimeID)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	col := database.DB.Collection("seats")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var docs []interface{}
	for _, row := range body.Rows {
		for i := 1; i <= body.SeatsPerRow; i++ {
			seat := models.Seat{
				ID:         primitive.NewObjectID(),
				ShowtimeID: showtimeID,
				Row:        row,
				Number:     i,
				Label:      fmt.Sprintf("%s%d", row, i),
				Status:     "available",
			}
			docs = append(docs, seat)
		}
	}

	if _, err := col.InsertMany(ctx, docs); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "seats generated",
		"total":   len(docs),
	})
}