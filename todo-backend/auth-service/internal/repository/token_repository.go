package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/todo-app/auth-service/internal/model"
	"gorm.io/gorm"
)

type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	FindRefreshToken(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	DeleteAllUserTokens(ctx context.Context, userID uuid.UUID) error
}

type tokenRepository struct{ db *gorm.DB }

func NewTokenRepository(db *gorm.DB) TokenRepository { return &tokenRepository{db} }

func (r *tokenRepository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	token := &model.RefreshToken{UserID: userID, TokenHash: tokenHash, ExpiresAt: expiresAt}
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *tokenRepository) FindRefreshToken(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	var t model.RefreshToken
	err := r.db.WithContext(ctx).Where("token_hash = ? AND expires_at > NOW()", tokenHash).First(&t).Error
	return &t, err
}

func (r *tokenRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	return r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).Delete(&model.RefreshToken{}).Error
}

func (r *tokenRepository) DeleteAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.RefreshToken{}).Error
}
