package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/todo-app/auth-service/internal/model"
	"github.com/todo-app/auth-service/internal/service"
)

// --- mocks ---

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) CreateUser(ctx context.Context, user *model.User) error {
	return m.Called(ctx, user).Error(0)
}
func (m *mockUserRepo) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *mockUserRepo) FindUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *mockUserRepo) UpdateMFASecret(ctx context.Context, id uuid.UUID, secret string, enabled bool) error {
	return m.Called(ctx, id, secret, enabled).Error(0)
}
func (m *mockUserRepo) SoftDeleteUser(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type mockTokenRepo struct{ mock.Mock }

func (m *mockTokenRepo) SaveRefreshToken(ctx context.Context, userID uuid.UUID, hash string, exp time.Time) error {
	return m.Called(ctx, userID, hash, exp).Error(0)
}
func (m *mockTokenRepo) FindRefreshToken(ctx context.Context, hash string) (*model.RefreshToken, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RefreshToken), args.Error(1)
}
func (m *mockTokenRepo) DeleteRefreshToken(ctx context.Context, hash string) error {
	return m.Called(ctx, hash).Error(0)
}
func (m *mockTokenRepo) DeleteAllUserTokens(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// --- helpers ---

func newTestService(t *testing.T) (service.AuthService, *mockUserRepo, *mockTokenRepo) {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	ur := &mockUserRepo{}
	tr := &mockTokenRepo{}
	svc := service.NewAuthService(ur, tr, rdb, "test-secret-32-chars-minimum-ok!")
	return svc, ur, tr
}

// --- tests ---

func TestRegister_DuplicateEmail(t *testing.T) {
	svc, ur, _ := newTestService(t)
	ur.On("FindUserByEmail", mock.Anything, "dup@example.com").Return(&model.User{}, nil)

	_, err := svc.Register(context.Background(), "dup@example.com", "password123")
	assert.ErrorIs(t, err, service.ErrUserExists)
}

func TestLogin_InvalidPassword(t *testing.T) {
	svc, ur, _ := newTestService(t)
	hash, _ := argon2id.CreateHash("correct", argon2id.DefaultParams)
	ur.On("FindUserByEmail", mock.Anything, "u@example.com").Return(&model.User{
		ID: uuid.New(), Email: "u@example.com", PasswordHash: hash,
	}, nil)

	_, err := svc.Login(context.Background(), "u@example.com", "wrong", "")
	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestLogin_AccountLocked(t *testing.T) {
	svc, ur, _ := newTestService(t)
	// Simulate 5 prior failures by calling login 5 times with wrong password
	hash, _ := argon2id.CreateHash("correct", argon2id.DefaultParams)
	ur.On("FindUserByEmail", mock.Anything, "locked@example.com").Return(&model.User{
		ID: uuid.New(), Email: "locked@example.com", PasswordHash: hash,
	}, nil)
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		svc.Login(ctx, "locked@example.com", "wrong", "")
	}
	_, err := svc.Login(ctx, "locked@example.com", "wrong", "")
	assert.ErrorIs(t, err, service.ErrAccountLocked)
}
