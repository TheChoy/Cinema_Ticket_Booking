package config

import (
    "log"
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    Port            string
    MongoURI        string
    MongoDB         string
    RabbitMQ        string
    RedisAddr       string
    JWTSecret       string
    FirebaseProject string
	FirebaseCredentials string
}

var C Config

func Load() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    C = Config{
        Port:            getEnv("PORT", "8081"),
        MongoURI:        getEnv("MONGODB_URI", ""),
        MongoDB:         getEnv("MONGODB_DB", "cinema"),
        RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
        RabbitMQ:        getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
        JWTSecret:       getEnv("JWT_SECRET", ""),
        FirebaseProject: getEnv("FIREBASE_PROJECT_ID", ""),
		FirebaseCredentials: getEnv("FIREBASE_CREDENTIALS", ""),
    }
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}