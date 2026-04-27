package service

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/rs/zerolog/log"
)

// Hub maintains active WebSocket connections keyed by userID.
type Hub struct {
	mu    sync.RWMutex
	conns map[string]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{conns: make(map[string]*websocket.Conn)}
}

func (h *Hub) Register(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	h.conns[userID] = conn
	h.mu.Unlock()
}

func (h *Hub) Unregister(userID string) {
	h.mu.Lock()
	delete(h.conns, userID)
	h.mu.Unlock()
}

// Send pushes a JSON message to the user's active connection.
// Returns false if the user has no active connection.
func (h *Hub) Send(userID string, payload any) bool {
	h.mu.RLock()
	conn, ok := h.conns[userID]
	h.mu.RUnlock()
	if !ok {
		return false
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return false
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("ws write failed")
		h.Unregister(userID)
		return false
	}
	return true
}
