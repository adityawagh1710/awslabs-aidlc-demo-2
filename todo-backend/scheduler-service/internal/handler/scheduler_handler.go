package handler

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/todo-app/scheduler-service/internal/service"
)

var validate = validator.New()

type SchedulerHandler struct{ svc service.SchedulerService }

func NewSchedulerHandler(svc service.SchedulerService) *SchedulerHandler {
	return &SchedulerHandler{svc}
}

type createReminderRequest struct {
	TodoID uuid.UUID `json:"todo_id" validate:"required"`
	UserID uuid.UUID `json:"user_id" validate:"required"`
	FireAt time.Time `json:"fire_at" validate:"required"`
}

type setRecurrenceRequest struct {
	CronExpression string `json:"cron_expression" validate:"required"`
}

func (h *SchedulerHandler) CreateReminder(c *fiber.Ctx) error {
	var req createReminderRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	if err := validate.Struct(req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	r, err := h.svc.ScheduleReminder(c.Context(), req.TodoID, req.UserID, req.FireAt)
	if err != nil {
		if err == service.ErrReminderLimit {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "reminder limit reached")
		}
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(r)
}

func (h *SchedulerHandler) DeleteReminder(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}
	return h.svc.CancelReminder(c.Context(), id)
}

func (h *SchedulerHandler) SetRecurrence(c *fiber.Ctx) error {
	todoID, err := uuid.Parse(c.Params("todoId"))
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}
	var req setRecurrenceRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	if err := validate.Struct(req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	if err := h.svc.SetRecurrence(c.Context(), todoID, req.CronExpression); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *SchedulerHandler) CompleteTodo(c *fiber.Ctx) error {
	todoID, err := uuid.Parse(c.Params("todoId"))
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}
	userID, _ := uuid.Parse(c.Query("user_id"))
	return h.svc.HandleTodoCompletion(c.Context(), todoID, userID)
}

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok", "service": "scheduler-service"})
}
