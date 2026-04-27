package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/todo-app/todo-service/internal/model"
	"gorm.io/gorm"
)

type TodoFilter struct {
	Status   model.TodoStatus
	Priority model.TodoPriority
	TagID    uuid.UUID
}

type TodoRepository interface {
	Insert(ctx context.Context, todo *model.Todo) error
	FindByUser(ctx context.Context, userID uuid.UUID, f TodoFilter) ([]model.Todo, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Todo, error)
	Update(ctx context.Context, todo *model.Todo) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type todoRepository struct{ db *gorm.DB }

func NewTodoRepository(db *gorm.DB) TodoRepository { return &todoRepository{db} }

func (r *todoRepository) Insert(ctx context.Context, todo *model.Todo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

func (r *todoRepository) FindByUser(ctx context.Context, userID uuid.UUID, f TodoFilter) ([]model.Todo, error) {
	q := r.db.WithContext(ctx).Preload("Tags").Where("user_id = ? AND deleted_at IS NULL", userID)
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Priority != "" {
		q = q.Where("priority = ?", f.Priority)
	}
	if f.TagID != uuid.Nil {
		q = q.Joins("JOIN todo_tags ON todo_tags.todo_id = todos.id").Where("todo_tags.tag_id = ?", f.TagID)
	}
	var todos []model.Todo
	return todos, q.Find(&todos).Error
}

func (r *todoRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Todo, error) {
	var todo model.Todo
	err := r.db.WithContext(ctx).Preload("Tags").Where("id = ? AND deleted_at IS NULL", id).First(&todo).Error
	return &todo, err
}

func (r *todoRepository) Update(ctx context.Context, todo *model.Todo) error {
	return r.db.WithContext(ctx).Save(todo).Error
}

func (r *todoRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.Todo{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}
