package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID         	  primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID     	  primitive.ObjectID   `bson:"user_id" json:"user_id"`
	ShowtimeID 	  primitive.ObjectID   `bson:"showtime_id" json:"showtime_id"`
	SeatIDs    	  []primitive.ObjectID `bson:"seat_ids" json:"seat_ids"`
	Status     	  string               `bson:"status" json:"status"`
	TotalPrice 	  float64              `bson:"total_price" json:"total_price"`
	CreatedAt  	  time.Time            `bson:"created_at" json:"created_at"`
	BookingNumber string 			   `bson:"booking_number" json:"booking_number"`
	PaidAt     	  *time.Time           `bson:"paid_at,omitempty" json:"paid_at,omitempty"`
}