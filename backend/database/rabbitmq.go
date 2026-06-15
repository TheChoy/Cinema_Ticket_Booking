package database

import (
	"encoding/json"
	"log"
	"time"

	"github.com/TheChoy/Cinema_Ticket_Booking/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitConn *amqp.Connection
var RabbitCh *amqp.Channel

const EventLogQueue = "event_logs"

func ConnectRabbitMQ() {
	var conn *amqp.Connection
	var err error

	// Retry สูงสุด 10 ครั้ง ครั้งละ 3 วินาที
	for i := 1; i <= 10; i++ {
		conn, err = amqp.Dial(config.C.RabbitMQ)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ connect error (attempt %d/10): %v", i, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("RabbitMQ failed after 10 attempts: %v", err)
	}
	RabbitConn = conn

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("RabbitMQ channel error: %v", err)
	}
	RabbitCh = ch

	_, err = ch.QueueDeclare(
		EventLogQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("RabbitMQ queue declare error: %v", err)
	}

	log.Println("Connected to RabbitMQ")
}

func PublishEventLog(data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Println("PublishEventLog marshal error:", err)
		return
	}

	err = RabbitCh.Publish(
		"",
		EventLogQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Println("PublishEventLog error:", err)
	}
}