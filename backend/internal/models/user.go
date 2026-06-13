package models

import (
    // "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UID      string             `bson:"uid" json:"uid"`
    Email    string             `bson:"email" json:"email"`
    Role     string             `bson:"role" json:"role"` // "user" | "admin"
    Password string             `bson:"-" json:"-"`
}

