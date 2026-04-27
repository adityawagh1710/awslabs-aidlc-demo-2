package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/todo-app/todo-service/internal/model"
	"github.com/todo-app/todo-service/internal/repository"
	"github.com/todo-app/todo-service/internal/service"
)

// --- mocks ---

type mockTodoRepo struct{ mock.Mock }

func (m *mockTodoRepo) Insert(ctx context.Context, t *model.Todo) error {
	return m.Called(ctx, t).Error(0)
}
func (m *mockTodoRepo) FindByUser(ctx context.Context, uid uuid.UUID, f repository.TodoFilter) ([]model.Todo, error) {
	args := m.Called(ctx, uid, f)
	return args.Get(0).([]model.Todo), args.Error(1)
}
func (m *mockTodoRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Todo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Todo), args.Error(1)
}
func (m *mockTodoRepo) Update(ctx context.Context, t *model.Todo) error {
	return m.Called(ctx, t).Error(0)
}
func (m *mockTodoRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type mockTagRepo struct{ mock.Mock }

func (m *mockTagRepo) Insert(ctx context.Context, t *model.Tag) error {
	return m.Called(ctx, t).Error(0)
}
func (m *mockTagRepo) FindByUser(ctx context.Context, uid uuid.UUID) ([]model.Tag, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]model.Tag), args.Error(1)
}
func (m *mockTagRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Tag, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}
func (m *mockTagRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type mockOutboxRepo struct{ mock.Mock }

func (m *mockOutboxRepo) Insert(ctx context.Context, e *model.SearchOutbox) error {
	return m.Called(ctx, e).Error(0)
}
func (m *mockOutboxRepo) FindUnprocessed(ctx context.Context, limit int) ([]model.SearchOutbox, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]model.SearchOutbox), args.Error(1)
}
func (m *mockOutboxRepo) MarkProcessed(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// --- tests ---

func newSvc(t *testing.T) (service.TodoService, *mockTodoRepo, *mockTagRepo, *mockOutboxRepo) {
	t.Helper()
	tr := &mockTodoRepo{}
	tgr := &mockTagRepo{}
	ob := &mockOutboxRepo{}
	svc := service.NewTodoService(tr, tgr, ob, nil, "", "")
	return svc, tr, tgr, ob
}

func TestCreate_Success(t *testing.T) {
	svc, tr, _, ob := newSvc(t)
	userID := uuid.New()
	tr.On("Insert", mock.Anything, mock.AnythingOfType("*model.Todo")).Return(nil)
	ob.On("Insert", mock.Anything, mock.Anything).Return(nil)

	todo, err := svc.Create(context.Background(), userID, service.CreateTodoInput{Title: "Test"})
	require.NoError(t, err)
	assert.Equal(t, "Test", todo.Title)
	assert.Equal(t, model.StatusPending, todo.Status)
}

func TestUpdate_InvalidTransition(t *testing.T) {
	svc, tr, _, _ := newSvc(t)
	userID := uuid.New()
	todoID := uuid.New()
	existing := &model.Todo{ID: todoID, UserID: userID, Status: model.StatusPending}
	tr.On("FindByID", mock.Anything, todoID).Return(existing, nil)

	done := model.StatusDone
	_, err := svc.Update(context.Background(), userID, todoID, service.UpdateTodoInput{Status: &done})
	assert.ErrorIs(t, err, service.ErrInvalidTransition)
}

func TestGet_Forbidden(t *testing.T) {
	svc, tr, _, _ := newSvc(t)
	todoID := uuid.New()
	tr.On("FindByID", mock.Anything, todoID).Return(&model.Todo{ID: todoID, UserID: uuid.New()}, nil)

	_, err := svc.Get(context.Background(), uuid.New(), todoID)
	assert.ErrorIs(t, err, service.ErrForbidden)
}
