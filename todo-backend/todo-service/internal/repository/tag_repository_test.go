package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/todo-app/todo-service/internal/model"
	"github.com/todo-app/todo-service/internal/repository"
	"github.com/todo-app/todo-service/internal/testutil"
)

func TestTagRepository(t *testing.T) {
	ctx := context.Background()
	pgc, db := testutil.StartPostgres(t, ctx)
	defer pgc.Terminate(ctx)

	repo := repository.NewTagRepository(db)
	userID := uuid.New()

	t.Run("Insert and FindByUser", func(t *testing.T) {
		tag := &model.Tag{UserID: userID, Name: "work"}
		require.NoError(t, repo.Insert(ctx, tag))
		assert.NotEmpty(t, tag.ID)

		tags, err := repo.FindByUser(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, tags, 1)
	})

	t.Run("Delete", func(t *testing.T) {
		tag := &model.Tag{UserID: userID, Name: "temp"}
		require.NoError(t, repo.Insert(ctx, tag))
		require.NoError(t, repo.Delete(ctx, tag.ID))

		tags, err := repo.FindByUser(ctx, userID)
		require.NoError(t, err)
		for _, tg := range tags {
			assert.NotEqual(t, tag.ID, tg.ID)
		}
	})
}
