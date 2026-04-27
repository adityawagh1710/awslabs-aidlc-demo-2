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

func TestTodoRepository(t *testing.T) {
	ctx := context.Background()
	pgc, db := testutil.StartPostgres(t, ctx)
	defer pgc.Terminate(ctx)

	repo := repository.NewTodoRepository(db)
	userID := uuid.New()

	t.Run("Insert and FindByUser", func(t *testing.T) {
		todo := &model.Todo{UserID: userID, Title: "Test todo", Status: model.StatusPending, Priority: model.PriorityMedium}
		require.NoError(t, repo.Insert(ctx, todo))
		assert.NotEmpty(t, todo.ID)

		todos, err := repo.FindByUser(ctx, userID, repository.TodoFilter{})
		require.NoError(t, err)
		assert.Len(t, todos, 1)
	})

	t.Run("FindByID", func(t *testing.T) {
		todo := &model.Todo{UserID: userID, Title: "Find me", Status: model.StatusPending, Priority: model.PriorityLow}
		require.NoError(t, repo.Insert(ctx, todo))

		found, err := repo.FindByID(ctx, todo.ID)
		require.NoError(t, err)
		assert.Equal(t, todo.ID, found.ID)
	})

	t.Run("Update", func(t *testing.T) {
		todo := &model.Todo{UserID: userID, Title: "Update me", Status: model.StatusPending, Priority: model.PriorityMedium}
		require.NoError(t, repo.Insert(ctx, todo))
		todo.Status = model.StatusInProgress
		require.NoError(t, repo.Update(ctx, todo))

		found, err := repo.FindByID(ctx, todo.ID)
		require.NoError(t, err)
		assert.Equal(t, model.StatusInProgress, found.Status)
	})

	t.Run("SoftDelete excludes from queries", func(t *testing.T) {
		todo := &model.Todo{UserID: userID, Title: "Delete me", Status: model.StatusPending, Priority: model.PriorityHigh}
		require.NoError(t, repo.Insert(ctx, todo))
		require.NoError(t, repo.SoftDelete(ctx, todo.ID))

		_, err := repo.FindByID(ctx, todo.ID)
		assert.Error(t, err)
	})

	t.Run("FilterByStatus", func(t *testing.T) {
		uid := uuid.New()
		require.NoError(t, repo.Insert(ctx, &model.Todo{UserID: uid, Title: "A", Status: model.StatusPending, Priority: model.PriorityLow}))
		require.NoError(t, repo.Insert(ctx, &model.Todo{UserID: uid, Title: "B", Status: model.StatusDone, Priority: model.PriorityLow}))

		todos, err := repo.FindByUser(ctx, uid, repository.TodoFilter{Status: model.StatusPending})
		require.NoError(t, err)
		assert.Len(t, todos, 1)
		assert.Equal(t, "A", todos[0].Title)
	})
}
