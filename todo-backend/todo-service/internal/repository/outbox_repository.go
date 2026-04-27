package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/todo-app/todo-service/internal/model"
	"gorm.io/gorm"
)

type OutboxRepository interface {
	Insert(ctx context.Context, event *model.SearchOutbox) error
	FindUnprocessed(ctx context.Context, limit int) ([]model.SearchOutbox, error)
	MarkProcessed(ctx context.Context, id uuid.UUID) error
}

type outboxRepository struct{ db *gorm.DB }

func NewOutboxRepository(db *gorm.DB) OutboxRepository { return &outboxRepository{db} }

func (r *outboxRepository) Insert(ctx context.Context, event *model.SearchOutbox) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *outboxRepository) FindUnprocessed(ctx context.Context, limit int) ([]model.SearchOutbox, error) {
	var events []model.SearchOutbox
	err := r.db.WithContext(ctx).
		Where("processed = false").
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

func (r *outboxRepository) MarkProcessed(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.SearchOutbox{}).
		Where("id = ?", id).
		Update("processed", true).Error
}
