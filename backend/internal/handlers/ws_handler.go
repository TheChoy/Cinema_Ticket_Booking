package handlers

import (
	"github.com/gofiber/websocket/v2"

	"github.com/TheChoy/Cinema_Ticket_Booking/internal/ws"
)

func WSSeatStatus(c *websocket.Conn) {
	showtimeID := c.Params("showtime_id")

	ws.H.Join(showtimeID, c)
	defer ws.H.Leave(showtimeID, c)

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}