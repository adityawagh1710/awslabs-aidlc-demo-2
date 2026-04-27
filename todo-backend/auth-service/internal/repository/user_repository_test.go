package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/todo-app/auth-service/internal/model"
	"github.com/todo-app/auth-service/internal/repository"
	"github.com/todo-app/auth-service/internal/testutil"
)

func TestUserRepository(t *testing.T) {
	ctx := context.Background()
	pgc, db := testutil.StartPostgres(t, ctx)
	defer pgc.Terminate(ctx)

	repo := repository.NewUserRepository(db)

	t.Run("CreateUser and FindByEmail", func(t *testing.T) {
		user := &model.User{Email: "test@example.com", PasswordHash: "hash"}
		require.NoError(t, repo.CreateUser(ctx, user))
		assert.NotEmpty(t, user.ID)

		found, err := repo.FindUserByEmail(ctx, "test@example.com")
		require.NoError(t, err)
		assert.Equal(t, user.ID, found.ID)
	})

	t.Run("FindByEmail returns error for unknown email", func(t *testing.T) {
		_, err := repo.FindUserByEmail(ctx, "nobody@example.com")
		assert.Error(t, err)
	})

	t.Run("UpdateMFASecret", func(t *testing.T) {
		user := &model.User{Email: "mfa@example.com", PasswordHash: "hash"}
		require.NoError(t, repo.CreateUser(ctx, user))

		require.NoError(t, repo.UpdateMFASecret(ctx, user.ID, "TOTP_SECRET", true))
		found, err := repo.FindUserByID(ctx, user.ID)
		require.NoError(t, err)
		assert.True(t, found.MFAEnabled)
		assert.Equal(t, "TOTP_SECRET", *found.MFASecret)
	})

	t.Run("SoftDeleteUser excludes from queries", func(t *testing.T) {
		user := &model.User{Email: "delete@example.com", PasswordHash: "hash"}
		require.NoError(t, repo.CreateUser(ctx, user))
		require.NoError(t, repo.SoftDeleteUser(ctx, user.ID))

		_, err := repo.FindUserByEmail(ctx, "delete@example.com")
		assert.Error(t, err)
	})
}

// Ensure postgres import used
var _ = (*postgres.PostgresContainer)(nil)
