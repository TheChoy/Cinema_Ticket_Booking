package main

import (
	"log"
	"time"
  "os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"

  "github.com/gofiber/swagger"
  _"github.com/TheChoy/Cinema_Ticket_Booking/docs"
)

// @title Book API
// @description This is a sample server for a book API.
// @version 1.0
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// Book struct to hold book data
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  app := fiber.New()
  
  app.Get("swagger/*", swagger.HandlerDefault)

  books = append(books, Book{ID : 1, Title: "Spai", Author: "TheChoy"})
  books = append(books, Book{ID : 2, Title: "Mod", Author: "TheChoy"})

  app.Post("/login", login)

  app.Use(checkMiddleware)

  app.Use(jwtware.New(jwtware.Config{
    SigningKey: []byte(os.Getenv("JWT_SECRET")),
  }))

  app.Get("/books",getBooks )
  app.Get("/books/:id", getBook)
  app.Post("/books", createBook)
  app.Put("/books/:id", updateBook)
  app.Delete("/books/:id", deleteBook)

  app.Listen(":8081")
}

type User struct {
  Email    string  `json:"email"`
  Password string  `json:"password"`
}
var memberUser = User{
  Email:    "user@example.com",
  Password: "password123",
}

func login(c * fiber.Ctx) error  {
  user := new(User)
  if err := c.BodyParser(user); err != nil {
    return c.Status(fiber.StatusBadRequest).SendString((err.Error()))
  }

  if user.Email != memberUser.Email || user.Password != memberUser.Password{
    return  fiber.ErrUnauthorized
  }

  // Create token
    token := jwt.New(jwt.SigningMethodHS256)

    // Set claims
    claims := token.Claims.(jwt.MapClaims)
    claims["email"] = user.Email
    claims["role"] = "admin"
    claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

    // Generate encoded token
    t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
      return c.SendStatus(fiber.StatusInternalServerError)
    }

  return c.JSON(fiber.Map{
    "masseage" : "Login Success",
    "token" : t,
  })
}

func checkMiddleware(c * fiber.Ctx) error  {
  user := c.Locals("user").(*jwt.Token)
  claims := user.Claims.(jwt.MapClaims)

  if claims["role"] != "admin" {
    return fiber.ErrUnauthorized
  }

  return c.Next()
}