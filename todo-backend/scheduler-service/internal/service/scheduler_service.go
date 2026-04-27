package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/sony/gobreaker"
	"github.com/todo-app/scheduler-service/internal/model"
	"github.com/todo-app/scheduler-service/internal/repository"
)

var (
	ErrReminderLimit = errors.New("reminder limit reached")
	ErrNotFound      = errors.New("not found")
)

type SchedulerService interface {
	ScheduleReminder(ctx context.Context, todoID, userID uuid.UUID, fireAt time.Time) (*model.Reminder, error)
	CancelReminder(ctx context.Context, reminderID uuid.UUID) error
	SetRecurrence(ctx context.Context, todoID uuid.UUID, cronExpr string) error
	HandleTodoCompletion(ctx context.Context, todoID, userID uuid.UUID) error
}

type schedulerService struct {
	reminderRepo   repository.ReminderRepository
	recurrenceRepo repository.RecurrenceRepository
	notifClient    *resty.Client
	todoClient     *resty.Client
	notifURL       string
	todoURL        string
}

func NewSchedulerService(
	rr repository.ReminderRepository,
	rcr repository.RecurrenceRepository,
	notifURL, todoURL string,
) SchedulerService {
	cb := func(name string) *resty.Client {
		cbk := gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    name,
			Timeout: 30 * time.Second,
			ReadyToTrip: func(c gobreaker.Counts) bool {
				return c.ConsecutiveFailures >= 5
			},
		})
		client := resty.New().SetRetryCount(3).SetRetryWaitTime(500 * time.Millisecond)
		client.OnBeforeRequest(func(_ *resty.Client, _ *resty.Request) error {
			_, err := cbk.Execute(func() (interface{}, error) { return nil, nil })
			return err
		})
		return client
	}
	return &schedulerService{rr, rcr, cb("notification"), cb("todo"), notifURL, todoURL}
}

func (s *schedulerService) ScheduleReminder(ctx context.Context, todoID, userID uuid.UUID, fireAt time.Time) (*model.Reminder, error) {
	count, err := s.reminderRepo.CountByTodo(ctx, todoID)
	if err != nil {
		return nil, err
	}
	if count >= 10 {
		return nil, ErrReminderLimit
	}
	r := &model.Reminder{TodoID: todoID, UserID: userID, FireAt: fireAt}
	return r, s.reminderRepo.Insert(ctx, r)
}

func (s *schedulerService) CancelReminder(ctx context.Context, reminderID uuid.UUID) error {
	return s.reminderRepo.Delete(ctx, reminderID)
}

func (s *schedulerService) SetRecurrence(ctx context.Context, todoID uuid.UUID, cronExpr string) error {
	next, err := NextOccurrence(cronExpr, time.Now())
	if err != nil {
		return err
	}
	rc := &model.RecurrenceConfig{TodoID: todoID, CronExpression: cronExpr, NextOccurrence: next}
	return s.recurrenceRepo.Upsert(ctx, rc)
}

func (s *schedulerService) HandleTodoCompletion(ctx context.Context, todoID, userID uuid.UUID) error {
	rc, err := s.recurrenceRepo.FindByTodo(ctx, todoID)
	if err != nil {
		return nil // no recurrence configured — not an error
	}
	next, err := NextOccurrence(rc.CronExpression, time.Now())
	if err != nil {
		return err
	}
	// Create next occurrence via todo-service
	_, err = s.todoClient.R().SetContext(ctx).
		SetBody(map[string]interface{}{"todo_id": todoID, "next_occurrence": next}).
		Post(s.todoURL + "/internal/todos/recur")
	if err != nil {
		log.Error().Err(err).Msg("failed to create recurrence todo")
	}
	// Delete old reminders
	return s.reminderRepo.DeleteByTodo(ctx, todoID)
}

// RunScheduler polls for due reminders every 30s and fires notification events.
func RunScheduler(ctx context.Context, rr repository.ReminderRepository, notifURL string) {
	client := resty.New().SetRetryCount(2)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fireReminders(ctx, rr, client, notifURL)
		}
	}
}

func fireReminders(ctx context.Context, rr repository.ReminderRepository, client *resty.Client, notifURL string) {
	due, err := rr.FindDue(ctx, time.Now())
	if err != nil {
		log.Error().Err(err).Msg("scheduler: fetch due reminders failed")
		return
	}
	for _, r := range due {
		if err := rr.MarkFired(ctx, r.ID); err != nil {
			continue
		}
		_, err := client.R().SetContext(ctx).
			SetBody(map[string]interface{}{
				"user_id": r.UserID,
				"todo_id": r.TodoID,
				"message": "Reminder for your todo",
			}).
			Post(notifURL + "/internal/events")
		if err != nil {
			log.Error().Err(err).Str("reminder_id", r.ID.String()).Msg("scheduler: notify failed")
		}
	}
}
