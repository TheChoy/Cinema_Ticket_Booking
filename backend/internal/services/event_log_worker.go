package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/TheChoy/Cinema_Ticket_Booking/database"
	"github.com/TheChoy/Cinema_Ticket_Booking/internal/models"
)

func StartEventLogWorker() {
	msgs, err := database.RabbitCh.Consume(
		database.EventLogQueue,
		"",
		true,  // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("StartEventLogWorker consume error:", err)
		return
	}

	go func() {
		for msg := range msgs {
			var eventLog models.EventLog
			if err := json.Unmarshal(msg.Body, &eventLog); err != nil {
				log.Println("EventLog unmarshal error:", err)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			col := database.DB.Collection("event_logs")
			if _, err := col.InsertOne(ctx, eventLog); err != nil {
				log.Println("EventLog insert error:", err)
			}
			cancel()
		}
	}()

	log.Println("Event Log Worker started")
}