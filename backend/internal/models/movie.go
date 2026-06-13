package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	PosterURL   string             `bson:"poster_url" json:"poster_url"`
	Duration    int                `bson:"duration" json:"duration"` // นาที
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	Status 		string			   `bson:"status" json:"status"`  //เผื่ออยากทำระบบให้ขึ้นหน้าคัมมิ่งซูน หรือกำลังฉาย
}