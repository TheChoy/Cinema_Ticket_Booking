package models

import (
    // "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
<<<<<<< HEAD
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UID      string             `bson:"uid" json:"uid"`
    Email    string             `bson:"email" json:"email"`
    Role     string             `bson:"role" json:"role"` // "user" | "admin"
    Password string             `bson:"-" json:"-"`
=======
    ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Email string             `bson:"email" json:"email"`
    Role  string             `bson:"role" json:"role"` // "user" | "admin"
    Password string           `json:"password"`
>>>>>>> 5497d4e04c73e5c7c543b1d5d5260635bdab344f
}

