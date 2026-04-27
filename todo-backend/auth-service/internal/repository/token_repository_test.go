package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/todo-app/auth-service/internal/model"
	"github.com/todo-app/auth-service/internal/repository"
	"github.com/todo-app/auth-service/internal/testutil"
)

func TestTokenRepository(t *testing.T) {
	ctx := context.Background()
	pgc, db := testutil.StartPostgres(t, ctx)
	defer pgc.Terminate(ctx)

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	user := &model.User{Email: "token@example.com", PasswordHash: "hash"}
	require.NoError(t, userRepo.CreateUser(ctx, user))

	t.Run("SaveAndFind", func(t *testing.T) {
		exp := time.Now().Add(30 * 24 * time.Hour)
		require.NoError(t, tokenRepo.SaveRefreshToken(ctx, user.ID, "hash1", exp))

		found, err := tokenRepo.FindRefreshToken(ctx, "hash1")
		require.NoError(t, err)
		assert.Equal(t, user.ID, found.UserID)
	})

	t.Run("DeleteRefreshToken", func(t *testing.T) {
		exp := time.Now().Add(30 * 24 * time.Hour)
		require.NoError(t, tokenRepo.SaveRefreshToken(ctx, user.ID, "hash2", exp))
		require.NoError(t, tokenRepo.DeleteRefreshToken(ctx, "hash2"))

		_, err := tokenRepo.FindRefreshToken(ctx, "hash2")
		assert.Error(t, err)
	})

	t.Run("ExpiredTokenNotFound", func(t *testing.T) {
		exp := time.Now().Add(-1 * time.Hour)
		require.NoError(t, tokenRepo.SaveRefreshToken(ctx, user.ID, "expired", exp))

		_, err := tokenRepo.FindRefreshToken(ctx, "expired")
		assert.Error(t, err)
	})
}
