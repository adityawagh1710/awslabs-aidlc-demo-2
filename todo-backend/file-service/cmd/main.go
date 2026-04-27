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

	"github.com/todo-app/file-service/internal/handler"
	"github.com/todo-app/file-service/internal/middleware"
	"github.com/todo-app/file-service/internal/repository"
	"github.com/todo-app/file-service/internal/service"
)

func main() {
	level, _ := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	zerolog.SetGlobalLevel(level)
	log.Logger = log.With().Str("service", "file-service").Logger()

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

	opts, _ := redis.ParseURL(os.Getenv("REDIS_URL"))
	rdb := redis.NewClient(opts)

	fileRepo := repository.NewFileRepository(db)
	fileSvc := service.NewFileService(fileRepo, os.Getenv("FILE_STORAGE_PATH"))
	fileH := handler.NewFileHandler(fileSvc)
	healthH := handler.NewHealthHandler()

	app := fiber.New(fiber.Config{
		ErrorHandler:  middleware.ErrorHandler,
		BodyLimit:     11 * 1024 * 1024, // 11MB to allow rejection at service layer
	})
	app.Use(middleware.Recover())
	jwtMW := middleware.JWTAuth(os.Getenv("JWT_SECRET"), rdb)

	app.Get("/health", healthH.Health)
	files := app.Group("/files", jwtMW)
	files.Post("/", fileH.Upload)
	files.Get("/:id", fileH.Download)
	files.Delete("/:id", fileH.Delete)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-quit
		shutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = app.ShutdownWithContext(shutCtx)
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
		_ = rdb.Close()
	}()

	log.Info().Msg("file-service listening on :3002")
	if err := app.Listen(":3002"); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}
