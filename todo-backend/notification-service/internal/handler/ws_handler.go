package handler

import (
	"context"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/todo-app/notification-service/internal/service"
)

type WSHandler struct {
	svc    service.NotificationService
	hub    *service.Hub
	secret string
}

func NewWSHandler(svc service.NotificationService, hub *service.Hub, secret string) *WSHandler {
	return &WSHandler{svc, hub, secret}
}

// Upgrade upgrades HTTP to WebSocket after JWT validation.
func (h *WSHandler) Upgrade(c *fiber.Ctx) error {
	tokenStr := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if tokenStr == "" {
		tokenStr = c.Query("token") // allow token as query param for WS
	}
	claims := &struct {
		UserID string `json:"user_id"`
		jwt.RegisteredClaims
	}{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(h.secret), nil
	})
	if err != nil || !token.Valid {
		return fiber.ErrUnauthorized
	}
	c.Locals("userID", claims.UserID)
	return websocket.New(h.handleWS)(c)
}

func (h *WSHandler) handleWS(c *websocket.Conn) {
	userID := c.Locals("userID").(string)
	uid, err := uuid.Parse(userID)
	if err != nil {
		return
	}

	h.hub.Register(userID, c)
	defer h.hub.Unregister(userID)

	// Push pending notifications on connect (NOTIF-03)
	pending, err := h.svc.GetPending(context.Background(), uid)
	if err == nil {
		for _, n := range pending {
			h.hub.Send(userID, n)
		}
	}

	// Keep connection alive — read loop (client pings or closes)
	for {
		if _, _, err := c.ReadMessage(); err != nil {
			log.Debug().Str("user_id", userID).Msg("ws connection closed")
			break
		}
	}
}
