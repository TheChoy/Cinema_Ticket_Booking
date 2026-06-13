package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Email string             `bson:"email" json:"email"`
    Role  string             `bson:"role" json:"role"` // "user" | "admin"
}

type Seat struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    ShowtimeID string             `bson:"showtime_id" json:"showtime_id"`
    SeatNumber string             `bson:"seat_number" json:"seat_number"`
    Status     string             `bson:"status" json:"status"` // AVAILABLE | LOCKED | BOOKED
}

type Booking struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID     string             `bson:"user_id" json:"user_id"`
    SeatID     string             `bson:"seat_id" json:"seat_id"`
    ShowtimeID string             `bson:"showtime_id" json:"showtime_id"`
    Status     string             `bson:"status" json:"status"` // PENDING | CONFIRMED | EXPIRED
    CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

type AuditLog struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Event     string             `bson:"event" json:"event"`
    UserID    string             `bson:"user_id" json:"user_id"`
    SeatID    string             `bson:"seat_id" json:"seat_id"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}