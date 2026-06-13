package middleware

import (
<<<<<<< HEAD
	"context"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/TheChoy/Cinema_Ticket_Booking/database"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/models"
)
var firebaseAuth *auth.Client

func InitFirebase() error {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return err
	}
	firebaseAuth, err = app.Auth(context.Background())
	return err
}

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		return fiber.ErrUnauthorized
	}
	idToken := authHeader[7:]

	token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	c.Locals("uid", token.UID)
	c.Locals("email", token.Claims["email"])
	return c.Next()
}

func AdminOnly(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)

	col := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	if err := col.FindOne(ctx, bson.M{"uid": uid}).Decode(&user); err != nil {
		return fiber.ErrUnauthorized
	}

	if user.Role != "admin" {
		return fiber.ErrForbidden
	}

=======
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func CheckMiddleware(c *fiber.Ctx) error {

	user := c.Locals("user")

	if user == nil {
		return fiber.ErrUnauthorized
	}

	token := user.(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)

	if claims["role"] != "admin" {
		return fiber.ErrUnauthorized
	}

>>>>>>> 5497d4e04c73e5c7c543b1d5d5260635bdab344f
	return c.Next()
}