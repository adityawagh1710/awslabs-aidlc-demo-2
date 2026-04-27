package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/todo-app/auth-service/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateMFASecret(ctx context.Context, userID uuid.UUID, secret string, enabled bool) error
	SoftDeleteUser(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository { return &userRepository{db} }

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&u).Error
	return &u, err
}

func (r *userRepository) FindUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&u).Error
	return &u, err
}

func (r *userRepository) UpdateMFASecret(ctx context.Context, userID uuid.UUID, secret string, enabled bool) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{"mfa_secret": secret, "mfa_enabled": enabled}).Error
}

func (r *userRepository) SoftDeleteUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{"is_active": false, "deleted_at": gorm.Expr("NOW()")}).Error
}
