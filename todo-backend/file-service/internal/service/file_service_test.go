package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/todo-app/file-service/internal/model"
	"github.com/todo-app/file-service/internal/service"
)

type mockFileRepo struct{ mock.Mock }

func (m *mockFileRepo) Insert(ctx context.Context, f *model.FileAttachment) error {
	return m.Called(ctx, f).Error(0)
}
func (m *mockFileRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.FileAttachment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.FileAttachment), args.Error(1)
}
func (m *mockFileRepo) FindByTodo(ctx context.Context, id uuid.UUID) ([]model.FileAttachment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]model.FileAttachment), args.Error(1)
}
func (m *mockFileRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockFileRepo) CountByTodo(ctx context.Context, id uuid.UUID) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func TestUpload_InvalidMIME(t *testing.T) {
	repo := &mockFileRepo{}
	svc := service.NewFileService(repo, t.TempDir())
	_, err := svc.Upload(context.Background(), uuid.New(), uuid.New(), "file.exe", "application/x-msdownload", []byte("data"))
	assert.ErrorIs(t, err, service.ErrInvalidMIME)
}

func TestUpload_TooLarge(t *testing.T) {
	repo := &mockFileRepo{}
	svc := service.NewFileService(repo, t.TempDir())
	data := make([]byte, model.MaxFileSizeBytes+1)
	_, err := svc.Upload(context.Background(), uuid.New(), uuid.New(), "big.pdf", "application/pdf", data)
	assert.ErrorIs(t, err, service.ErrFileTooLarge)
}

func TestUpload_Success(t *testing.T) {
	repo := &mockFileRepo{}
	dir := t.TempDir()
	svc := service.NewFileService(repo, dir)

	userID, todoID := uuid.New(), uuid.New()
	repo.On("CountByTodo", mock.Anything, todoID).Return(int64(0), nil)
	repo.On("Insert", mock.Anything, mock.AnythingOfType("*model.FileAttachment")).Return(nil)

	f, err := svc.Upload(context.Background(), userID, todoID, "doc.pdf", "application/pdf", []byte("content"))
	require.NoError(t, err)
	assert.Equal(t, "doc.pdf", f.Filename)

	// Verify file written to disk
	_, statErr := os.Stat(f.StoragePath)
	assert.NoError(t, statErr)
}

func TestDelete_Forbidden(t *testing.T) {
	repo := &mockFileRepo{}
	svc := service.NewFileService(repo, t.TempDir())
	fileID := uuid.New()
	repo.On("FindByID", mock.Anything, fileID).Return(&model.FileAttachment{ID: fileID, UserID: uuid.New()}, nil)

	err := svc.Delete(context.Background(), uuid.New(), fileID)
	assert.ErrorIs(t, err, service.ErrForbidden)
}
