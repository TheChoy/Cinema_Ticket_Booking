package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Showtime struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MovieID   primitive.ObjectID `bson:"movie_id" json:"movie_id"`
	Room      string             `bson:"room" json:"room"`
	StartTime time.Time          `bson:"start_time" json:"start_time"`
	EndTime   time.Time          `bson:"end_time" json:"end_time"`
	SeatCount int                `bson:"seat_count" json:"seat_count"`
	Price 	  float64 			 `bson:"price" json:"price"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}