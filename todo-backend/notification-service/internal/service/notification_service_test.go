package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/todo-app/notification-service/internal/model"
	"github.com/todo-app/notification-service/internal/service"
)

type mockNotifRepo struct{ mock.Mock }

func (m *mockNotifRepo) Insert(ctx context.Context, n *model.Notification) error {
	return m.Called(ctx, n).Error(0)
}
func (m *mockNotifRepo) FindPendingByUser(ctx context.Context, id uuid.UUID) ([]model.Notification, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]model.Notification), args.Error(1)
}
func (m *mockNotifRepo) MarkDelivered(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func TestDeliver_UserOffline_StoresUndelivered(t *testing.T) {
	repo := &mockNotifRepo{}
	hub := service.NewHub() // no connections registered
	svc := service.NewNotificationService(repo, hub)

	userID := uuid.New()
	repo.On("Insert", mock.Anything, mock.MatchedBy(func(n *model.Notification) bool {
		return n.UserID == userID && !n.Delivered
	})).Return(nil)

	err := svc.Deliver(context.Background(), userID, nil, "Reminder!")
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestGetPending(t *testing.T) {
	repo := &mockNotifRepo{}
	hub := service.NewHub()
	svc := service.NewNotificationService(repo, hub)

	userID := uuid.New()
	expected := []model.Notification{{ID: uuid.New(), UserID: userID, Message: "test"}}
	repo.On("FindPendingByUser", mock.Anything, userID).Return(expected, nil)

	result, err := svc.GetPending(context.Background(), userID)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}
