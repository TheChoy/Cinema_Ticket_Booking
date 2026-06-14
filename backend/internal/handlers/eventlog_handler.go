package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/TheChoy/Cinema_Ticket_Booking/database"
)

func GetEventLogs(c *fiber.Ctx) error {
	col := database.DB.Collection("event_logs")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	matchStage := bson.M{}

	if event := c.Query("event"); event != "" {
		matchStage["event"] = event
	}

	if userID := c.Query("user_id"); userID != "" {
		id, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		matchStage["user_id"] = id
	}

	if bookingID := c.Query("booking_id"); bookingID != "" {
		id, err := primitive.ObjectIDFromHex(bookingID)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		matchStage["booking_id"] = id
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{Key: "$sort", Value: bson.M{"created_at": -1}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$user", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "bookings",
			"localField":   "booking_id",
			"foreignField": "_id",
			"as":           "booking",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$booking", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "seats",
			"localField":   "seat_ids",
			"foreignField": "_id",
			"as":           "seats",
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "showtimes",
			"localField":   "booking.showtime_id",
			"foreignField": "_id",
			"as":           "showtime",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$showtime", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "movies",
			"localField":   "showtime.movie_id",
			"foreignField": "_id",
			"as":           "movie",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$movie", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$project", Value: bson.M{
			"_id":            1,
			"event":          1,
			"message":        1,
			"created_at":     1,
			"user_email":     "$user.email",
			"booking_number": "$booking.booking_number",
			"movie_title":    "$movie.title",
			"seats":          "$seats.label",
		}}},
	}

	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(results)
}