package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/TheChoy/Cinema_Ticket_Booking/database"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/ws"
)

func StartBookingExpirer() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			expireBookings()
		}
	}()

	go func() {
		for {
			time.Sleep(10 * time.Second)
			broadcastCountdowns()
		}
	}()
}
func broadcastCountdowns() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bookingCol := database.DB.Collection("bookings")

	cursor, err := bookingCol.Find(ctx, bson.M{"status": "pending"})
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	type Booking struct {
		ID         primitive.ObjectID `bson:"_id"`
		ShowtimeID primitive.ObjectID `bson:"showtime_id"`
		CreatedAt  time.Time          `bson:"created_at"`
	}

	var bookings []Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return
	}

	for _, b := range bookings {
		expireAt := b.CreatedAt.Add(5 * time.Minute)
		remaining := int(time.Until(expireAt).Seconds())
		if remaining < 0 {
			remaining = 0
		}

		ws.H.Broadcast(b.ShowtimeID.Hex(), ws.Message{
			Type:             "countdown",
			BookingID:        b.ID.Hex(),
			RemainingSeconds: remaining,
		})
	}
}

func expireBookings() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bookingCol := database.DB.Collection("bookings")
	seatCol := database.DB.Collection("seats")

	expiredTime := time.Now().Add(-5 * time.Minute)
	cursor, err := bookingCol.Find(ctx, bson.M{
		"status":     "pending",
		"created_at": bson.M{"$lt": expiredTime},
	})
	if err != nil {
		log.Println("expireBookings error:", err)
		return
	}
	defer cursor.Close(ctx)

	type Booking struct {
		ID         interface{}   `bson:"_id"`
		ShowtimeID interface{}   `bson:"showtime_id"`
		SeatIDs    []interface{} `bson:"seat_ids"`
	}

	var bookings []Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return
	}

	for _, b := range bookings {
		seatCol.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": b.SeatIDs}}, bson.M{"$set": bson.M{
			"status":     "available",
			"booking_id": nil,
		}})

		bookingCol.UpdateOne(ctx, bson.M{"_id": b.ID}, bson.M{"$set": bson.M{"status": "cancelled"}})

		// Broadcast seat update
		showtimeID := b.ShowtimeID.(primitive.ObjectID).Hex()
		for _, seatID := range b.SeatIDs {
			ws.H.Broadcast(showtimeID, ws.Message{
				Type:   "seat_update",
				SeatID: seatID.(primitive.ObjectID).Hex(),
				Status: "available",
			})
		}

		log.Printf("Booking %v expired and cancelled", b.ID)
	}
}