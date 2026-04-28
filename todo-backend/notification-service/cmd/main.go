package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/todo-app/notification-service/internal/handler"
	"github.com/todo-app/notification-service/internal/middleware"
	"github.com/todo-app/notification-service/internal/repository"
	"github.com/todo-app/notification-service/internal/service"
)

func main() {
	level, _ := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	zerolog.SetGlobalLevel(level)
	log.Logger = log.With().Str("service", "notification-service").Logger()

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("db connect failed")
	}
	m, err := migrate.New("file://migrations", func() string {
		if u := os.Getenv("MIGRATIONS_DB_URL"); u != "" {
			return u
		}
		return os.Getenv("DATABASE_URL")
	}())
	if err != nil {
		log.Fatal().Err(err).Msg("migrate init failed")
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("migrate up failed")
	}

	hub := service.NewHub()
	notifRepo := repository.NewNotificationRepository(db)
	notifSvc := service.NewNotificationService(notifRepo, hub)

	wsH := handler.NewWSHandler(notifSvc, hub, os.Getenv("JWT_SECRET"))
	eventH := handler.NewEventHandler(notifSvc)
	healthH := handler.NewHealthHandler()
	apiKey := middleware.InternalAPIKey(os.Getenv("INTERNAL_API_KEY"))

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Use(middleware.Recover())

	prom := fiberprometheus.New("notification-service")
	prom.RegisterAt(app, "/metrics")
	app.Use(prom.Middleware)

	app.Get("/health", healthH.Health)
	app.Get("/ws", wsH.Upgrade)
	app.Post("/internal/events", apiKey, eventH.Ingest)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-quit
		shutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = app.ShutdownWithContext(shutCtx)
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}()

	log.Info().Msg("notification-service listening on :3003")
	if err := app.Listen(":3003"); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}
