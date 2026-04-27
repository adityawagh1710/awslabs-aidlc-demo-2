package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
	"github.com/todo-app/auth-service/internal/middleware"
	"github.com/todo-app/auth-service/internal/model"
	"github.com/todo-app/auth-service/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountLocked      = errors.New("account locked")
	ErrMFARequired        = errors.New("mfa required")
	ErrInvalidMFA         = errors.New("invalid mfa code")
	ErrUserExists         = errors.New("email already registered")
)

const (
	accessTokenTTL  = time.Hour
	refreshTokenTTL = 30 * 24 * time.Hour
	maxLoginAttempts = 5
	lockoutDuration  = 15 * time.Minute
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthService interface {
	Register(ctx context.Context, email, password string) (*TokenPair, error)
	Login(ctx context.Context, email, password, mfaCode string) (*TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*TokenPair, error)
	Logout(ctx context.Context, userID, accessToken, refreshToken string) error
	EnrollMFA(ctx context.Context, userID string) (secret, qrURL string, err error)
	VerifyMFA(ctx context.Context, userID, code string) error
}

type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	rdb       *redis.Client
	jwtSecret string
}

func NewAuthService(u repository.UserRepository, t repository.TokenRepository, rdb *redis.Client, secret string) AuthService {
	return &authService{u, t, rdb, secret}
}

func (s *authService) Register(ctx context.Context, email, password string) (*TokenPair, error) {
	if _, err := s.userRepo.FindUserByEmail(ctx, email); err == nil {
		return nil, ErrUserExists
	}
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}
	user := &model.User{Email: email, PasswordHash: hash}
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return s.issueTokenPair(ctx, user.ID.String())
}

func (s *authService) Login(ctx context.Context, email, password, mfaCode string) (*TokenPair, error) {
	lockKey := fmt.Sprintf("lockout:%s", email)
	attempts, _ := s.rdb.Get(ctx, lockKey).Int()
	if attempts >= maxLoginAttempts {
		return nil, ErrAccountLocked
	}

	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		s.incrementAttempts(ctx, lockKey)
		return nil, ErrInvalidCredentials
	}

	match, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil || !match {
		s.incrementAttempts(ctx, lockKey)
		return nil, ErrInvalidCredentials
	}

	if user.MFAEnabled {
		if mfaCode == "" {
			return nil, ErrMFARequired
		}
		if user.MFASecret == nil || !totp.Validate(mfaCode, *user.MFASecret) {
			return nil, ErrInvalidMFA
		}
	}

	s.rdb.Del(ctx, lockKey)
	return s.issueTokenPair(ctx, user.ID.String())
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	h := hashToken(refreshToken)
	stored, err := s.tokenRepo.FindRefreshToken(ctx, h)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := s.tokenRepo.DeleteRefreshToken(ctx, h); err != nil {
		return nil, err
	}
	return s.issueTokenPair(ctx, stored.UserID.String())
}

func (s *authService) Logout(ctx context.Context, userID, accessToken, refreshToken string) error {
	// Blacklist access token for its remaining TTL
	claims := &middleware.Claims{}
	t, _ := jwt.ParseWithClaims(accessToken, claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if t != nil && t.Valid {
		ttl := time.Until(claims.ExpiresAt.Time)
		if ttl > 0 {
			s.rdb.Set(ctx, "blacklist:"+accessToken, 1, ttl)
		}
	}
	return s.tokenRepo.DeleteRefreshToken(ctx, hashToken(refreshToken))
}

func (s *authService) EnrollMFA(ctx context.Context, userID string) (string, string, error) {
	user, err := s.userRepo.FindUserByID(ctx, uuid.MustParse(userID))
	if err != nil {
		return "", "", err
	}
	key, err := totp.Generate(totp.GenerateOpts{Issuer: "TodoApp", AccountName: user.Email})
	if err != nil {
		return "", "", err
	}
	if err := s.userRepo.UpdateMFASecret(ctx, user.ID, key.Secret(), false); err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

func (s *authService) VerifyMFA(ctx context.Context, userID, code string) error {
	user, err := s.userRepo.FindUserByID(ctx, uuid.MustParse(userID))
	if err != nil {
		return err
	}
	if user.MFASecret == nil || !totp.Validate(code, *user.MFASecret) {
		return ErrInvalidMFA
	}
	return s.userRepo.UpdateMFASecret(ctx, user.ID, *user.MFASecret, true)
}

// issueTokenPair creates a new access + refresh token pair.
func (s *authService) issueTokenPair(ctx context.Context, userID string) (*TokenPair, error) {
	now := time.Now()
	claims := &middleware.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.New().String()
	if err := s.tokenRepo.SaveRefreshToken(ctx, uuid.MustParse(userID), hashToken(refreshToken), now.Add(refreshTokenTTL)); err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *authService) incrementAttempts(ctx context.Context, key string) {
	s.rdb.Incr(ctx, key)
	s.rdb.Expire(ctx, key, lockoutDuration)
}

func hashToken(token string) string {
	return HashToken(token)
}

func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
