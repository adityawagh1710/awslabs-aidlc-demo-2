package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/todo-app/todo-service/internal/model"
	"gorm.io/gorm"
)

type TagRepository interface {
	Insert(ctx context.Context, tag *model.Tag) error
	FindByUser(ctx context.Context, userID uuid.UUID) ([]model.Tag, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Tag, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type tagRepository struct{ db *gorm.DB }

func NewTagRepository(db *gorm.DB) TagRepository { return &tagRepository{db} }

func (r *tagRepository) Insert(ctx context.Context, tag *model.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *tagRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]model.Tag, error) {
	var tags []model.Tag
	return tags, r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&tags).Error
}

func (r *tagRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Tag, error) {
	var tag model.Tag
	return &tag, r.db.WithContext(ctx).First(&tag, "id = ?", id).Error
}

func (r *tagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Tag{}, "id = ?", id).Error
}
