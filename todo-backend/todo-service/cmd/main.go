package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/todo-app/todo-service/internal/handler"
	"github.com/todo-app/todo-service/internal/middleware"
	"github.com/todo-app/todo-service/internal/repository"
	"github.com/todo-app/todo-service/internal/service"
)

func main() {
	level, _ := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	zerolog.SetGlobalLevel(level)
	log.Logger = log.With().Str("service", "todo-service").Logger()

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

	es, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{os.Getenv("ELASTICSEARCH_URL")}})
	if err != nil {
		log.Fatal().Err(err).Msg("es connect failed")
	}

	todoRepo := repository.NewTodoRepository(db)
	tagRepo := repository.NewTagRepository(db)
	outboxRepo := repository.NewOutboxRepository(db)

	todoSvc := service.NewTodoService(todoRepo, tagRepo, outboxRepo, es,
		os.Getenv("FILE_SERVICE_URL"), os.Getenv("SCHEDULER_SERVICE_URL"))
	tagSvc := service.NewTagService(tagRepo)

	// Start outbox worker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go service.RunOutboxWorker(ctx, outboxRepo, es)

	todoH := handler.NewTodoHandler(todoSvc)
	tagH := handler.NewTagHandler(tagSvc)
	healthH := handler.NewHealthHandler()

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Use(middleware.RequestLogger())
	app.Use(middleware.Recover())

	jwtMW := middleware.JWTAuth(os.Getenv("JWT_SECRET"), rdb)
	app.Get("/health", healthH.Health)
	todos := app.Group("/todos", jwtMW)
	todos.Post("/", todoH.Create)
	todos.Get("/", todoH.List)
	todos.Get("/search", todoH.Search)
	todos.Get("/:id", todoH.Get)
	todos.Put("/:id", todoH.Update)
	todos.Delete("/:id", todoH.Delete)

	tags := app.Group("/tags", jwtMW)
	tags.Post("/", tagH.Create)
	tags.Get("/", tagH.List)
	tags.Delete("/:id", tagH.Delete)

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
		_ = rdb.Close()
	}()

	log.Info().Msg("todo-service listening on :3001")
	if err := app.Listen(":3001"); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}
