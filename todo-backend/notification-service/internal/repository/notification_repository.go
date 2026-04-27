package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/todo-app/notification-service/internal/model"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Insert(ctx context.Context, n *model.Notification) error
	FindPendingByUser(ctx context.Context, userID uuid.UUID) ([]model.Notification, error)
	MarkDelivered(ctx context.Context, id uuid.UUID) error
}

type notificationRepository struct{ db *gorm.DB }

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db}
}

func (r *notificationRepository) Insert(ctx context.Context, n *model.Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *notificationRepository) FindPendingByUser(ctx context.Context, userID uuid.UUID) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND delivered = false", userID).
		Order("created_at ASC").
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) MarkDelivered(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("id = ?", id).Update("delivered", true).Error
}
