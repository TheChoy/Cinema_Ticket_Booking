package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/TheChoy/Cinema_Ticket_Booking/database"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/models"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/ws"
)

func CreateBooking(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)

	body := struct {
		ShowtimeID string   `json:"showtime_id"`
		SeatIDs    []string `json:"seat_ids"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	showtimeID, err := primitive.ObjectIDFromHex(body.ShowtimeID)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var seatIDs []primitive.ObjectID
	for _, s := range body.SeatIDs {
		id, err := primitive.ObjectIDFromHex(s)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		seatIDs = append(seatIDs, id)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Redis Lock ที่นั่งทุกตัว
	for _, seatID := range body.SeatIDs {
		lockKey := "seat_lock:" + seatID
		ok, err := database.RDB.SetNX(ctx, lockKey, uid, 5*time.Minute).Result()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if !ok {
			// Publish event log: lock_failed
			seatID, _ := primitive.ObjectIDFromHex(seatID)
			database.PublishEventLog(models.EventLog{
				ID:        primitive.NewObjectID(),
				Event:     "lock_failed",
				SeatIDs:   []primitive.ObjectID{seatID},
				Message:   "Seat lock failed - already locked by another user",
				CreatedAt: time.Now(),
			})
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "seat " + lockKey + " is being booked by someone else"})
		}
	}

	// เช็คว่าที่นั่งว่างทั้งหมด
	seatCol := database.DB.Collection("seats")
	cursor, err := seatCol.Find(ctx, bson.M{
		"_id":    bson.M{"$in": seatIDs},
		"status": "available",
	})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	var availableSeats []models.Seat
	cursor.All(ctx, &availableSeats)

	if len(availableSeats) != len(seatIDs) {
		// ปลด lock ถ้าที่นั่งไม่ว่าง
		for _, seatID := range body.SeatIDs {
			database.RDB.Del(ctx, "seat_lock:"+seatID)
		}
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "some seats are not available"})
	}

	// ดึง showtime เพื่อคำนวณราคา
	showtimeCol := database.DB.Collection("showtimes")
	var showtime models.Showtime
	if err := showtimeCol.FindOne(ctx, bson.M{"_id": showtimeID}).Decode(&showtime); err != nil {
		return fiber.ErrNotFound
	}

	totalPrice := showtime.Price * float64(len(seatIDs))

	// ดึง user ID จาก MongoDB
	userCol := database.DB.Collection("users")
	var user models.User
	if err := userCol.FindOne(ctx, bson.M{"uid": uid}).Decode(&user); err != nil {
		return fiber.ErrNotFound
	}

	booking := models.Booking{
		ID:            primitive.NewObjectID(),
		UserID:        user.ID,
		ShowtimeID:    showtimeID,
		SeatIDs:       seatIDs,
		Status:        "pending",
		TotalPrice:    totalPrice,
		CreatedAt:     time.Now(),
		BookingNumber: fmt.Sprintf("BK%d", time.Now().UnixNano()),
	}

	bookingCol := database.DB.Collection("bookings")
	if _, err := bookingCol.InsertOne(ctx, booking); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	
	// Publish event log: booking_success
	database.PublishEventLog(models.EventLog{
		ID:        primitive.NewObjectID(),
		Event:     "booking_success",
		UserID:    &user.ID,
		BookingID: &booking.ID,
		SeatIDs:   seatIDs,
		Message:   fmt.Sprintf("User %s booked %d seat(s)", user.Email, len(seatIDs)),
		CreatedAt: time.Now(),
	})

	// อัพเดท status ที่นั่งเป็น locked
	seatCol.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": seatIDs}}, bson.M{"$set": bson.M{
		"status":     "locked",
		"booking_id": booking.ID,
	}})

	// Broadcast seat update
	for _, seatID := range body.SeatIDs {
		ws.H.Broadcast(body.ShowtimeID, ws.Message{
			Type:   "seat_update",
			SeatID: seatID,
			Status: "locked",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(booking)
}

func GetMyBookings(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCol := database.DB.Collection("users")
	var user models.User
	if err := userCol.FindOne(ctx, bson.M{"uid": uid}).Decode(&user); err != nil {
		return fiber.ErrNotFound
	}

	bookingCol := database.DB.Collection("bookings")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": user.ID}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "showtimes",
			"localField":   "showtime_id",
			"foreignField": "_id",
			"as":           "showtime",
		}}},
		{{Key: "$unwind", Value: "$showtime"}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "movies",
			"localField":   "showtime.movie_id",
			"foreignField": "_id",
			"as":           "movie",
		}}},
		{{Key: "$unwind", Value: "$movie"}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "seats",
			"localField":   "seat_ids",
			"foreignField": "_id",
			"as":           "seats",
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":            1,
			"booking_number": 1,
			"status":         1,
			"total_price":    1,
			"paid_at":        1,
			"movie_title":    "$movie.title",
			"poster_url":     "$movie.poster_url",
			"showtime":       "$showtime.start_time",
			"seats":          "$seats.label",
		}}},
	}

	cursor, err := bookingCol.Aggregate(ctx, pipeline)
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

func GetBookingByID(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bookingCol := database.DB.Collection("bookings")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": id}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "showtimes",
			"localField":   "showtime_id",
			"foreignField": "_id",
			"as":           "showtime",
		}}},
		{{Key: "$unwind", Value: "$showtime"}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "movies",
			"localField":   "showtime.movie_id",
			"foreignField": "_id",
			"as":           "movie",
		}}},
		{{Key: "$unwind", Value: "$movie"}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "seats",
			"localField":   "seat_ids",
			"foreignField": "_id",
			"as":           "seats",
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":            1,
			"booking_number": 1,
			"status":         1,
			"total_price":    1,
			"created_at":     1,
			"showtime_id":    1,
			"movie_title":    "$movie.title",
			"poster_url":     "$movie.poster_url",
			"room":           "$showtime.room",
			"start_time":     "$showtime.start_time",
			"end_time":       "$showtime.end_time",
			"seats":          "$seats.label",
		}}},
	}

	cursor, err := bookingCol.Aggregate(ctx, pipeline)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if len(results) == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(results[0])
}

func PayBooking(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bookingCol := database.DB.Collection("bookings")

	var booking models.Booking
	if err := bookingCol.FindOne(ctx, bson.M{"_id": id, "status": "pending"}).Decode(&booking); err != nil {
		return fiber.ErrNotFound
	}

	now := time.Now()
	if _, err := bookingCol.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"status":  "paid",
		"paid_at": now,
	}}); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// เปลี่ยน status ที่นั่งจาก locked เป็น booked
	seatCol := database.DB.Collection("seats")
	seatCol.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": booking.SeatIDs}}, bson.M{"$set": bson.M{
		"status": "booked",
	}})

	// Broadcast seat update
	for _, seatID := range booking.SeatIDs {
		ws.H.Broadcast(booking.ShowtimeID.Hex(), ws.Message{
			Type:   "seat_update",
			SeatID: seatID.Hex(),
			Status: "booked",
		})
	}
	// ปลด Redis Lock
	for _, seatID := range booking.SeatIDs {
		database.RDB.Del(ctx, "seat_lock:"+seatID.Hex())
	}

	return c.JSON(fiber.Map{"message": "payment successful"})
}

func GetAllBookings(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bookingCol := database.DB.Collection("bookings")
	cursor, err := bookingCol.Find(ctx, bson.M{})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(bookings)
}

func UpdateBooking(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	body := struct {
		Status string `json:"status"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bookingCol := database.DB.Collection("bookings")
	if _, err := bookingCol.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": body.Status}}); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"message": "booking updated"})
}

func CancelBooking(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bookingCol := database.DB.Collection("bookings")

	// ดึง booking ก่อนเพื่อเอา seat IDs
	var booking models.Booking
	if err := bookingCol.FindOne(ctx, bson.M{"_id": id}).Decode(&booking); err != nil {
		return fiber.ErrNotFound
	}

	// คืน status ที่นั่งเป็น available
	seatCol := database.DB.Collection("seats")
	seatCol.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": booking.SeatIDs}}, bson.M{"$set": bson.M{
		"status":     "available",
		"booking_id": nil,
	}})

	// Broadcast seat update
	for _, seatID := range booking.SeatIDs {
		ws.H.Broadcast(booking.ShowtimeID.Hex(), ws.Message{
			Type:   "seat_update",
			SeatID: seatID.Hex(),
			Status: "available",
		})
	}

	// อัพเดท status booking เป็น cancelled
	bookingCol.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": "cancelled"}})

	// Publish event log: seat_release
	database.PublishEventLog(models.EventLog{
		ID:        primitive.NewObjectID(),
		Event:     "seat_release",
		BookingID: &booking.ID,
		SeatIDs:   booking.SeatIDs,
		Message:   "Booking cancelled by admin, seats released",
		CreatedAt: time.Now(),
	})
	
	// ปลด Redis Lock
	for _, seatID := range booking.SeatIDs {
		database.RDB.Del(ctx, "seat_lock:"+seatID.Hex())
	}

	return c.JSON(fiber.Map{"message": "booking cancelled"})
}