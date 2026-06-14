package handlers

import (
    "context"
    "fmt"
    "mime/multipart"
    "time"
	"io"
	"log"
	"os"

    // "cloud.google.com/go/storage"
    "github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    // "google.golang.org/api/option"

    // "github.com/TheChoy/Cinema_Ticket_Booking/config"
    "github.com/TheChoy/Cinema_Ticket_Booking/database"
    "github.com/TheChoy/Cinema_Ticket_Booking/internal/models"
)

func uploadPoster(file *multipart.FileHeader) (string, error) {
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	savePath := "./uploads/" + filename

	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	url := fmt.Sprintf("http://localhost:8081/uploads/%s", filename)
	return url, nil
}

func CreateMovie(c *fiber.Ctx) error {
	movie := &models.Movie{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Genre:       c.FormValue("genre"),
		Status:      c.FormValue("status"),
	}

	if duration := c.FormValue("duration"); duration != "" {
		fmt.Sscanf(duration, "%d", &movie.Duration)
	}

	// รับรูปถ้ามี
	if file, err := c.FormFile("poster"); err == nil {
		url, err := uploadPoster(file)
		if err != nil {
			log.Println("uploadPoster error:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		movie.PosterURL = url
	}

	movie.ID = primitive.NewObjectID()
	movie.CreatedAt = time.Now()

	col := database.DB.Collection("movies")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := col.InsertOne(ctx, movie); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(movie)
}

func GetMovies(c *fiber.Ctx) error {
	col := database.DB.Collection("movies")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}

	if search := c.Query("search"); search != "" {
		filter["title"] = bson.M{"$regex": search, "$options": "i"}
	}
	if genre := c.Query("genre"); genre != "" {
		filter["genre"] = genre
	}
	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer cursor.Close(ctx)

	var movies []models.Movie
	if err := cursor.All(ctx, &movies); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(movies)
}

func UpdateMovie(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	movie := new(models.Movie)
	if err := c.BodyParser(movie); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	col := database.DB.Collection("movies")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{
		"title":       movie.Title,
		"description": movie.Description,
		"poster_url":  movie.PosterURL,
		"duration":    movie.Duration,
		"status":      movie.Status,
		"genre":       movie.Genre,
	}}

	if _, err := col.UpdateOne(ctx, bson.M{"_id": id}, update); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"message": "updated"})
}

func DeleteMovie(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	col := database.DB.Collection("movies")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := col.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"message": "deleted"})
}