package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/todo-app/scheduler-service/internal/model"
	"github.com/todo-app/scheduler-service/internal/repository"
	"github.com/todo-app/scheduler-service/internal/testutil"
)

func TestReminderRepository(t *testing.T) {
	ctx := context.Background()
	_, db := testutil.StartPostgres(t, ctx)
	repo := repository.NewReminderRepository(db)

	userID := uuid.New()
	todoID := uuid.New()

	t.Run("Insert and FindDue", func(t *testing.T) {
		// Insert a reminder due in the past (should be found)
		past := &model.Reminder{
			TodoID: todoID, UserID: userID,
			FireAt: time.Now().Add(-1 * time.Minute),
		}
		require.NoError(t, repo.Insert(ctx, past))

		// Insert a reminder due in the future (should NOT be found)
		future := &model.Reminder{
			TodoID: uuid.New(), UserID: userID,
			FireAt: time.Now().Add(1 * time.Hour),
		}
		require.NoError(t, repo.Insert(ctx, future))

		due, err := repo.FindDue(ctx, time.Now())
		require.NoError(t, err)
		ids := make([]uuid.UUID, len(due))
		for i, r := range due { ids[i] = r.ID }
		assert.Contains(t, ids, past.ID)
		assert.NotContains(t, ids, future.ID)
	})

	t.Run("MarkFired excludes from FindDue", func(t *testing.T) {
		r := &model.Reminder{
			TodoID: uuid.New(), UserID: userID,
			FireAt: time.Now().Add(-30 * time.Second),
		}
		require.NoError(t, repo.Insert(ctx, r))

		// Found before marking
		due, _ := repo.FindDue(ctx, time.Now())
		ids := make([]uuid.UUID, len(due))
		for i, d := range due { ids[i] = d.ID }
		assert.Contains(t, ids, r.ID)

		require.NoError(t, repo.MarkFired(ctx, r.ID))

		// Not found after marking
		due, _ = repo.FindDue(ctx, time.Now())
		ids = make([]uuid.UUID, len(due))
		for i, d := range due { ids[i] = d.ID }
		assert.NotContains(t, ids, r.ID)
	})

	t.Run("Delete", func(t *testing.T) {
		r := &model.Reminder{
			TodoID: uuid.New(), UserID: userID,
			FireAt: time.Now().Add(-1 * time.Second),
		}
		require.NoError(t, repo.Insert(ctx, r))
		require.NoError(t, repo.Delete(ctx, r.ID))

		due, _ := repo.FindDue(ctx, time.Now())
		for _, d := range due {
			assert.NotEqual(t, r.ID, d.ID)
		}
	})

	t.Run("CountByTodo", func(t *testing.T) {
		tid := uuid.New()
		for i := 0; i < 3; i++ {
			r := &model.Reminder{TodoID: tid, UserID: userID, FireAt: time.Now().Add(time.Hour)}
			require.NoError(t, repo.Insert(ctx, r))
		}
		count, err := repo.CountByTodo(ctx, tid)
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})

	t.Run("DeleteByTodo removes all reminders for a todo", func(t *testing.T) {
		tid := uuid.New()
		for i := 0; i < 2; i++ {
			r := &model.Reminder{TodoID: tid, UserID: userID, FireAt: time.Now().Add(-time.Second)}
			require.NoError(t, repo.Insert(ctx, r))
		}
		require.NoError(t, repo.DeleteByTodo(ctx, tid))

		count, err := repo.CountByTodo(ctx, tid)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestRecurrenceRepository(t *testing.T) {
	ctx := context.Background()
	_, db := testutil.StartPostgres(t, ctx)
	repo := repository.NewRecurrenceRepository(db)

	t.Run("Upsert and FindByTodo", func(t *testing.T) {
		todoID := uuid.New()
		rc := &model.RecurrenceConfig{
			TodoID:         todoID,
			CronExpression: "0 9 * * 1", // every Monday 9am
			NextOccurrence: time.Now().Add(7 * 24 * time.Hour),
		}
		require.NoError(t, repo.Upsert(ctx, rc))
		assert.NotEqual(t, uuid.Nil, rc.ID)

		got, err := repo.FindByTodo(ctx, todoID)
		require.NoError(t, err)
		assert.Equal(t, "0 9 * * 1", got.CronExpression)
	})

	t.Run("Upsert updates existing config", func(t *testing.T) {
		todoID := uuid.New()
		rc := &model.RecurrenceConfig{
			TodoID:         todoID,
			CronExpression: "0 9 * * 1",
			NextOccurrence: time.Now().Add(7 * 24 * time.Hour),
		}
		require.NoError(t, repo.Upsert(ctx, rc))

		rc.CronExpression = "0 10 * * 2"
		require.NoError(t, repo.Upsert(ctx, rc))

		got, err := repo.FindByTodo(ctx, todoID)
		require.NoError(t, err)
		assert.Equal(t, "0 10 * * 2", got.CronExpression)
	})

	t.Run("Delete", func(t *testing.T) {
		todoID := uuid.New()
		rc := &model.RecurrenceConfig{
			TodoID:         todoID,
			CronExpression: "@daily",
			NextOccurrence: time.Now().Add(24 * time.Hour),
		}
		require.NoError(t, repo.Upsert(ctx, rc))
		require.NoError(t, repo.Delete(ctx, todoID))

		_, err := repo.FindByTodo(ctx, todoID)
		assert.Error(t, err)
	})
}
