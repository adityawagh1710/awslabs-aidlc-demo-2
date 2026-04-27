package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/todo-app/auth-service/internal/handler"
	"github.com/todo-app/auth-service/internal/middleware"
	"github.com/todo-app/auth-service/internal/repository"
	"github.com/todo-app/auth-service/internal/service"
)

func main() {
	// Logger
	level, _ := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	zerolog.SetGlobalLevel(level)
	log.Logger = log.With().Str("service", "auth-service").Logger()

	// DB
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// Migrations
	migrationsURL := os.Getenv("MIGRATIONS_DB_URL")
	if migrationsURL == "" {
		migrationsURL = os.Getenv("DATABASE_URL")
	}
	m, err := migrate.New("file://migrations", migrationsURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init migrations")
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	// Redis
	opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse redis url")
	}
	rdb := redis.NewClient(opts)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}

	// Wire dependencies
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	authSvc := service.NewAuthService(userRepo, tokenRepo, rdb, os.Getenv("JWT_SECRET"))
	authHandler := handler.NewAuthHandler(authSvc)
	healthHandler := handler.NewHealthHandler()

	// Fiber app — trust X-Real-IP / X-Forwarded-For set by nginx / Traefik
	// so c.IP() returns the real client address, not the proxy container IP.
	app := fiber.New(fiber.Config{
		ErrorHandler:          middleware.ErrorHandler,
		ProxyHeader:           "X-Real-IP",
		EnableTrustedProxyCheck: true,
		TrustedProxies:        []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"},
	})
	app.Use(middleware.RequestLogger())
	app.Use(middleware.RateLimiter(rdb))
	app.Use(middleware.Recover())

	// Routes
	app.Get("/health", healthHandler.Health)
	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", middleware.JWTAuth(os.Getenv("JWT_SECRET"), rdb), authHandler.Logout)
	auth.Post("/mfa/enroll", middleware.JWTAuth(os.Getenv("JWT_SECRET"), rdb), authHandler.MFAEnroll)
	auth.Post("/mfa/verify", middleware.JWTAuth(os.Getenv("JWT_SECRET"), rdb), authHandler.MFAVerify)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-quit
		log.Info().Msg("shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = app.ShutdownWithContext(ctx)
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
		_ = rdb.Close()
	}()

	log.Info().Msg("auth-service listening on :3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}
