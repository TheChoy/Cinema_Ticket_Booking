package database

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/TheChoy/Cinema_Ticket_Booking/config"
)

var RDB *redis.Client

func ConnectRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: config.C.RedisAddr,
	})

	ctx := context.Background()
	if _, err := RDB.Ping(ctx).Result(); err != nil {
		log.Fatalf("Redis connect error: %v", err)
	}

	log.Println("Connected to Redis")
}