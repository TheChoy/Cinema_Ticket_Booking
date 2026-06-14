package ws

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Message struct {
	Type   string `json:"type"`   // "seat_update" | "countdown"
	SeatID string `json:"seat_id,omitempty"`
	Status string `json:"status,omitempty"`
	BookingID string `json:"booking_id,omitempty"`
	RemainingSeconds int `json:"remaining_seconds,omitempty"`
}

type Hub struct {
	mu      sync.Mutex
	rooms   map[string]map[*websocket.Conn]bool // showtime_id -> connections
}

var H = &Hub{
	rooms: make(map[string]map[*websocket.Conn]bool),
}

func (h *Hub) Join(showtimeID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[showtimeID] == nil {
		h.rooms[showtimeID] = make(map[*websocket.Conn]bool)
	}
	h.rooms[showtimeID][conn] = true
}

func (h *Hub) Leave(showtimeID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.rooms[showtimeID], conn)
}

func (h *Hub) Broadcast(showtimeID string, msg Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	for conn := range h.rooms[showtimeID] {
		conn.WriteMessage(websocket.TextMessage, data)
	}
}