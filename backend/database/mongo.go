package database

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "github.com/TheChoy/Cinema_Ticket_Booking/config"
)
var DB *mongo.Database

func ConnectMongo() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.C.MongoURI))
    if err != nil {
        log.Fatalf("MongoDB connect error: %v", err)
    }

    if err = client.Ping(ctx, nil); err != nil {
        log.Fatalf("MongoDB ping error: %v", err)
    }

    DB = client.Database(config.C.MongoDB)
    log.Printf("Connected to MongoDB: %s", config.C.MongoDB)
}