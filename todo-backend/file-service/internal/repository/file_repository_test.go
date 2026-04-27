package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/todo-app/file-service/internal/model"
	"github.com/todo-app/file-service/internal/repository"
	"github.com/todo-app/file-service/internal/testutil"
)

func TestFileRepository(t *testing.T) {
	ctx := context.Background()
	_, db := testutil.StartPostgres(t, ctx)
	repo := repository.NewFileRepository(db)

	userID := uuid.New()
	todoID := uuid.New()

	// Write a temp file so StoragePath refers to something real
	tmp, err := os.CreateTemp(t.TempDir(), "attachment-*.pdf")
	require.NoError(t, err)
	tmp.WriteString("pdf content")
	tmp.Close()

	t.Run("Insert and FindByID", func(t *testing.T) {
		f := &model.FileAttachment{
			TodoID:      todoID,
			UserID:      userID,
			Filename:    "test.pdf",
			StoragePath: tmp.Name(),
			MimeType:    "application/pdf",
			SizeBytes:   11,
		}
		require.NoError(t, repo.Insert(ctx, f))
		assert.NotEqual(t, uuid.Nil, f.ID)

		got, err := repo.FindByID(ctx, f.ID)
		require.NoError(t, err)
		assert.Equal(t, "test.pdf", got.Filename)
		assert.Equal(t, userID, got.UserID)
	})

	t.Run("FindByTodo returns all attachments", func(t *testing.T) {
		f := &model.FileAttachment{
			TodoID: todoID, UserID: userID,
			Filename: "doc.pdf", StoragePath: tmp.Name(),
			MimeType: "application/pdf", SizeBytes: 5,
		}
		require.NoError(t, repo.Insert(ctx, f))

		files, err := repo.FindByTodo(ctx, todoID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(files), 2)
	})

	t.Run("CountByTodo", func(t *testing.T) {
		count, err := repo.CountByTodo(ctx, todoID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(2))
	})

	t.Run("Delete removes attachment", func(t *testing.T) {
		f := &model.FileAttachment{
			TodoID: uuid.New(), UserID: userID,
			Filename: "todel.pdf", StoragePath: tmp.Name(),
			MimeType: "application/pdf", SizeBytes: 3,
		}
		require.NoError(t, repo.Insert(ctx, f))

		require.NoError(t, repo.Delete(ctx, f.ID))
		_, err := repo.FindByID(ctx, f.ID)
		assert.Error(t, err, "should be not found after delete")
	})
}
