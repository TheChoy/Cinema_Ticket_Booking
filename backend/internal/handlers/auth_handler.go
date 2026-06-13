package handlers

import (
<<<<<<< HEAD
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/TheChoy/Cinema_Ticket_Booking/database"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/models"
)

func Register(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	email := c.Locals("email").(string)

	col := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing models.User
	err := col.FindOne(ctx, bson.M{"uid": uid}).Decode(&existing)
	if err == nil {
		return c.JSON(existing)
	}
	if err != mongo.ErrNoDocuments {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	user := models.User{
		ID:    primitive.NewObjectID(),
		UID:   uid,
		Email: email,
		Role:  "user",
	}

	if _, err := col.InsertOne(ctx, user); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

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
=======
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"github.com/TheChoy/Cinema_Ticket_Booking/internal/models"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/services"
)

func Login(c *fiber.Ctx) error {

	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if user.Email != services.MemberUser.Email ||
		user.Password != services.MemberUser.Password {

		return fiber.ErrUnauthorized
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["role"] = "admin"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Login Success",
		"token":   t,
	})
>>>>>>> 5497d4e04c73e5c7c543b1d5d5260635bdab344f
}