package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventLog struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Event     string               `bson:"event" json:"event"`   
	UserID    *primitive.ObjectID  `bson:"user_id,omitempty" json:"user_id,omitempty"`
	BookingID *primitive.ObjectID  `bson:"booking_id,omitempty" json:"booking_id,omitempty"`
	SeatIDs   []primitive.ObjectID `bson:"seat_ids,omitempty" json:"seat_ids,omitempty"`
	Message   string               `bson:"message" json:"message"`       
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
}