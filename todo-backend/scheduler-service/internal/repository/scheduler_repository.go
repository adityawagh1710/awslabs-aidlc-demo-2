package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/todo-app/scheduler-service/internal/model"
	"gorm.io/gorm"
)

type ReminderRepository interface {
	Insert(ctx context.Context, r *model.Reminder) error
	FindDue(ctx context.Context, now time.Time) ([]model.Reminder, error)
	MarkFired(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByTodo(ctx context.Context, todoID uuid.UUID) error
	CountByTodo(ctx context.Context, todoID uuid.UUID) (int64, error)
}

type reminderRepository struct{ db *gorm.DB }

func NewReminderRepository(db *gorm.DB) ReminderRepository { return &reminderRepository{db} }

func (r *reminderRepository) Insert(ctx context.Context, rem *model.Reminder) error {
	return r.db.WithContext(ctx).Create(rem).Error
}
func (r *reminderRepository) FindDue(ctx context.Context, now time.Time) ([]model.Reminder, error) {
	var reminders []model.Reminder
	err := r.db.WithContext(ctx).Where("fire_at <= ? AND fired = false", now).Find(&reminders).Error
	return reminders, err
}
func (r *reminderRepository) MarkFired(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.Reminder{}).Where("id = ?", id).Update("fired", true).Error
}
func (r *reminderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Reminder{}, "id = ?", id).Error
}
func (r *reminderRepository) DeleteByTodo(ctx context.Context, todoID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Reminder{}, "todo_id = ?", todoID).Error
}
func (r *reminderRepository) CountByTodo(ctx context.Context, todoID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Reminder{}).Where("todo_id = ? AND fired = false", todoID).Count(&count).Error
	return count, err
}

type RecurrenceRepository interface {
	Upsert(ctx context.Context, rc *model.RecurrenceConfig) error
	FindByTodo(ctx context.Context, todoID uuid.UUID) (*model.RecurrenceConfig, error)
	Delete(ctx context.Context, todoID uuid.UUID) error
}

type recurrenceRepository struct{ db *gorm.DB }

func NewRecurrenceRepository(db *gorm.DB) RecurrenceRepository { return &recurrenceRepository{db} }

func (r *recurrenceRepository) Upsert(ctx context.Context, rc *model.RecurrenceConfig) error {
	return r.db.WithContext(ctx).Save(rc).Error
}
func (r *recurrenceRepository) FindByTodo(ctx context.Context, todoID uuid.UUID) (*model.RecurrenceConfig, error) {
	var rc model.RecurrenceConfig
	return &rc, r.db.WithContext(ctx).Where("todo_id = ?", todoID).First(&rc).Error
}
func (r *recurrenceRepository) Delete(ctx context.Context, todoID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.RecurrenceConfig{}, "todo_id = ?", todoID).Error
}
