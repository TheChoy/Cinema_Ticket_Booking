package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/TheChoy/Cinema_Ticket_Booking/database"
)

func StartBookingExpirer() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			expireBookings()
		}
	}()
}

func expireBookings() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bookingCol := database.DB.Collection("bookings")
	seatCol := database.DB.Collection("seats")

	// หา booking ที่ pending เกิน 5 นาที
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
		ID      interface{}   `bson:"_id"`
		SeatIDs []interface{} `bson:"seat_ids"`
	}

	var bookings []Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return
	}

	for _, b := range bookings {
		// คืน seat เป็น available
		seatCol.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": b.SeatIDs}}, bson.M{"$set": bson.M{
			"status":     "available",
			"booking_id": nil,
		}})

		// ยกเลิก booking
		bookingCol.UpdateOne(ctx, bson.M{"_id": b.ID}, bson.M{"$set": bson.M{"status": "cancelled"}})

		log.Printf("Booking %v expired and cancelled", b.ID)
	}
}