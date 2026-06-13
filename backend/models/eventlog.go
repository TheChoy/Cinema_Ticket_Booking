package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventLog struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Event     string                 `bson:"event" json:"event"`
	UserID    primitive.ObjectID     `bson:"user_id" json:"user_id"`
    SeatID    string                 `bson:"seat_id" json:"seat_id"`
	BookingID *primitive.ObjectID 	 `bson:"booking_id,omitempty" json:"booking_id,omitempty"`
	Payload	  map[string]interface{} `bson:"payload" json:"payload"`
	CreatedAt time.Time              `bson:"created_at" json:"created_at"`
}