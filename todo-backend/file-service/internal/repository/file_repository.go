package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/todo-app/file-service/internal/model"
	"gorm.io/gorm"
)

type FileRepository interface {
	Insert(ctx context.Context, f *model.FileAttachment) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.FileAttachment, error)
	FindByTodo(ctx context.Context, todoID uuid.UUID) ([]model.FileAttachment, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CountByTodo(ctx context.Context, todoID uuid.UUID) (int64, error)
}

type fileRepository struct{ db *gorm.DB }

func NewFileRepository(db *gorm.DB) FileRepository { return &fileRepository{db} }

func (r *fileRepository) Insert(ctx context.Context, f *model.FileAttachment) error {
	return r.db.WithContext(ctx).Create(f).Error
}
func (r *fileRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.FileAttachment, error) {
	var f model.FileAttachment
	return &f, r.db.WithContext(ctx).First(&f, "id = ?", id).Error
}
func (r *fileRepository) FindByTodo(ctx context.Context, todoID uuid.UUID) ([]model.FileAttachment, error) {
	var files []model.FileAttachment
	return files, r.db.WithContext(ctx).Where("todo_id = ?", todoID).Find(&files).Error
}
func (r *fileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.FileAttachment{}, "id = ?", id).Error
}
func (r *fileRepository) CountByTodo(ctx context.Context, todoID uuid.UUID) (int64, error) {
	var count int64
	return count, r.db.WithContext(ctx).Model(&model.FileAttachment{}).Where("todo_id = ?", todoID).Count(&count).Error
}
