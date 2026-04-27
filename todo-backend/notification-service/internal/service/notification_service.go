package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/todo-app/notification-service/internal/model"
	"github.com/todo-app/notification-service/internal/repository"
)

type NotificationService interface {
	Deliver(ctx context.Context, userID uuid.UUID, todoID *uuid.UUID, message string) error
	GetPending(ctx context.Context, userID uuid.UUID) ([]model.Notification, error)
}

type notificationService struct {
	repo repository.NotificationRepository
	hub  *Hub
}

func NewNotificationService(repo repository.NotificationRepository, hub *Hub) NotificationService {
	return &notificationService{repo, hub}
}

func (s *notificationService) Deliver(ctx context.Context, userID uuid.UUID, todoID *uuid.UUID, message string) error {
	n := &model.Notification{UserID: userID, TodoID: todoID, Message: message}

	// Try live WebSocket delivery first (NOTIF-01)
	if s.hub.Send(userID.String(), n) {
		n.Delivered = true
	}

	// Always persist — if delivered=false, will be sent on next connect (NOTIF-02)
	return s.repo.Insert(ctx, n)
}

func (s *notificationService) GetPending(ctx context.Context, userID uuid.UUID) ([]model.Notification, error) {
	return s.repo.FindPendingByUser(ctx, userID)
}
