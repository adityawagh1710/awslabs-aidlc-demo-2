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

	"github.com/todo-app/scheduler-service/internal/handler"
	"github.com/todo-app/scheduler-service/internal/middleware"
	"github.com/todo-app/scheduler-service/internal/repository"
	"github.com/todo-app/scheduler-service/internal/service"
)

func main() {
	level, _ := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	zerolog.SetGlobalLevel(level)
	log.Logger = log.With().Str("service", "scheduler-service").Logger()

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

	reminderRepo := repository.NewReminderRepository(db)
	recurrenceRepo := repository.NewRecurrenceRepository(db)
	svc := service.NewSchedulerService(reminderRepo, recurrenceRepo,
		os.Getenv("NOTIFICATION_SERVICE_URL"), os.Getenv("TODO_SERVICE_URL"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go service.RunScheduler(ctx, reminderRepo, os.Getenv("NOTIFICATION_SERVICE_URL"))

	h := handler.NewSchedulerHandler(svc)
	healthH := handler.NewHealthHandler()
	apiKey := middleware.InternalAPIKey(os.Getenv("INTERNAL_API_KEY"))

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Use(middleware.Recover())

	prom := fiberprometheus.New("scheduler-service")
	prom.RegisterAt(app, "/metrics")
	app.Use(prom.Middleware)

	app.Get("/health", healthH.Health)
	app.Post("/reminders", apiKey, h.CreateReminder)
	app.Delete("/reminders/:id", apiKey, h.DeleteReminder)
	app.Post("/todos/:todoId/recurrence", apiKey, h.SetRecurrence)
	app.Post("/todos/:todoId/complete", apiKey, h.CompleteTodo)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-quit
		cancel()
		shutCtx, shutCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutCancel()
		_ = app.ShutdownWithContext(shutCtx)
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}()

	log.Info().Msg("scheduler-service listening on :3004")
	if err := app.Listen(":3004"); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}
