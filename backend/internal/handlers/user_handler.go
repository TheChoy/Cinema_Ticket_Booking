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

func GetMe(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)

	col := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	if err := col.FindOne(ctx, bson.M{"uid": uid}).Decode(&user); err != nil {
		return fiber.ErrNotFound
	}

	return c.JSON(user)
}

func UpdateMe(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)

	body := struct {
		Name        string     `json:"name"`
		Phone       string     `json:"phone"`
		DateOfBirth *time.Time `json:"date_of_birth"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	col := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{
		"name":          body.Name,
		"phone":         body.Phone,
		"date_of_birth": body.DateOfBirth,
	}}

	if _, err := col.UpdateOne(ctx, bson.M{"uid": uid}, update); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"message": "updated"})
}

func GetUsers(c *fiber.Ctx) error {
	col := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(users)
}

func UpdateUserRole(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	body := struct {
		Role string `json:"role"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if body.Role != "user" && body.Role != "admin" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "role must be user or admin"})
	}

	col := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"role": body.Role}}); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"message": "role updated"})
}