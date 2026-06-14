package database

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitConn *amqp.Connection
var RabbitCh *amqp.Channel

const EventLogQueue = "event_logs"

func ConnectRabbitMQ() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("RabbitMQ connect error: %v", err)
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