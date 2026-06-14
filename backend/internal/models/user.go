package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UID         string             `bson:"uid" json:"uid"`
	Email       string             `bson:"email" json:"email"`
	Role        string             `bson:"role" json:"role"`
	Name        string             `bson:"name" json:"name"`
	Phone       string             `bson:"phone" json:"phone"`
	DateOfBirth *time.Time         `bson:"date_of_birth,omitempty" json:"date_of_birth,omitempty"`
	Password    string             `bson:"-" json:"-"`
}
