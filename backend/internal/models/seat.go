package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Seat struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	ShowtimeID  primitive.ObjectID  `bson:"showtime_id" json:"showtime_id"`
	Label       string              `bson:"label" json:"label"`
	Row         string              `bson:"row" json:"row"`
	Number      int                 `bson:"number" json:"number"`
	Status      string              `bson:"status" json:"status"`
	LockedBy    *primitive.ObjectID `bson:"locked_by,omitempty" json:"locked_by,omitempty"`
	LockedUntil *time.Time          `bson:"locked_until,omitempty" json:"locked_until,omitempty"`
	BookingID   *primitive.ObjectID `bson:"booking_id,omitempty" json:"booking_id,omitempty"`
}