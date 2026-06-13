package middleware

import (
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

	return c.Next()
}