package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/TheChoy/Cinema_Ticket_Booking/database"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/models"
)

func GetShowtimes(c *fiber.Ctx) error {
	col := database.DB.Collection("showtimes")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}

	if movieID := c.Query("movie_id"); movieID != "" {
		id, err := primitive.ObjectIDFromHex(movieID)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		filter["movie_id"] = id
	}

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer cursor.Close(ctx)

	var showtimes []models.Showtime
	if err := cursor.All(ctx, &showtimes); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(showtimes)
}

func CreateShowtime(c *fiber.Ctx) error {
	showtime := new(models.Showtime)
	if err := c.BodyParser(showtime); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	showtime.ID = primitive.NewObjectID()
	showtime.CreatedAt = time.Now()

	col := database.DB.Collection("showtimes")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := col.InsertOne(ctx, showtime); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(showtime)
}

func UpdateShowtime(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	showtime := new(models.Showtime)
	if err := c.BodyParser(showtime); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	col := database.DB.Collection("showtimes")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{
		"movie_id":   showtime.MovieID,
		"room":       showtime.Room,
		"start_time": showtime.StartTime,
		"end_time":   showtime.EndTime,
		"seat_count": showtime.SeatCount,
		"price":      showtime.Price,
	}}

	if _, err := col.UpdateOne(ctx, bson.M{"_id": id}, update); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"message": "updated"})
}

func DeleteShowtime(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	col := database.DB.Collection("showtimes")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := col.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"message": "deleted"})
}