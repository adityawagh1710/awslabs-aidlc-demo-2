package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/todo-app/notification-service/internal/service"
)

type EventHandler struct{ svc service.NotificationService }

func NewEventHandler(svc service.NotificationService) *EventHandler { return &EventHandler{svc} }

type ingestEventRequest struct {
	UserID  uuid.UUID  `json:"user_id"`
	TodoID  *uuid.UUID `json:"todo_id"`
	Message string     `json:"message"`
}

func (h *EventHandler) Ingest(c *fiber.Ctx) error {
	var req ingestEventRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	if req.UserID == uuid.Nil || req.Message == "" {
		return fiber.ErrUnprocessableEntity
	}
	if err := h.svc.Deliver(c.Context(), req.UserID, req.TodoID, req.Message); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusAccepted)
}

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok", "service": "notification-service"})
}
