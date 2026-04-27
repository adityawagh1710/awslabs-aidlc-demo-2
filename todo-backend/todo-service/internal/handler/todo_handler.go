package handler

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/todo-app/todo-service/internal/model"
	"github.com/todo-app/todo-service/internal/repository"
	"github.com/todo-app/todo-service/internal/service"
)

var validate = validator.New()

type TodoHandler struct{ svc service.TodoService }

func NewTodoHandler(svc service.TodoService) *TodoHandler { return &TodoHandler{svc} }

type createTodoRequest struct {
	Title       string             `json:"title" validate:"required,max=255"`
	Description string             `json:"description" validate:"max=5000"`
	Priority    model.TodoPriority `json:"priority"`
	DueDate     *time.Time         `json:"due_date"`
	TagIDs      []uuid.UUID        `json:"tag_ids"`
}

type updateTodoRequest struct {
	Title       *string            `json:"title" validate:"omitempty,max=255"`
	Description *string            `json:"description" validate:"omitempty,max=5000"`
	Status      *model.TodoStatus  `json:"status"`
	Priority    *model.TodoPriority `json:"priority"`
	DueDate     *time.Time         `json:"due_date"`
	TagIDs      []uuid.UUID        `json:"tag_ids"`
}

func (h *TodoHandler) Create(c *fiber.Ctx) error {
	var req createTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	if err := validate.Struct(req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	userID := mustUserID(c)
	if req.Priority == "" {
		req.Priority = model.PriorityMedium
	}
	todo, err := h.svc.Create(c.Context(), userID, service.CreateTodoInput{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		TagIDs:      req.TagIDs,
	})
	if err != nil {
		return mapErr(err)
	}
	return c.Status(fiber.StatusCreated).JSON(todo)
}

func (h *TodoHandler) List(c *fiber.Ctx) error {
	userID := mustUserID(c)
	f := repository.TodoFilter{
		Status:   model.TodoStatus(c.Query("status")),
		Priority: model.TodoPriority(c.Query("priority")),
	}
	if tagID := c.Query("tag_id"); tagID != "" {
		f.TagID, _ = uuid.Parse(tagID)
	}
	todos, err := h.svc.List(c.Context(), userID, f)
	if err != nil {
		return err
	}
	return c.JSON(todos)
}

func (h *TodoHandler) Get(c *fiber.Ctx) error {
	userID := mustUserID(c)
	todoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}
	todo, err := h.svc.Get(c.Context(), userID, todoID)
	if err != nil {
		return mapErr(err)
	}
	return c.JSON(todo)
}

func (h *TodoHandler) Update(c *fiber.Ctx) error {
	var req updateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	if err := validate.Struct(req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	userID := mustUserID(c)
	todoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}
	todo, err := h.svc.Update(c.Context(), userID, todoID, service.UpdateTodoInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		TagIDs:      req.TagIDs,
	})
	if err != nil {
		return mapErr(err)
	}
	return c.JSON(todo)
}

func (h *TodoHandler) Delete(c *fiber.Ctx) error {
	userID := mustUserID(c)
	todoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}
	if err := h.svc.Delete(c.Context(), userID, todoID); err != nil {
		return mapErr(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TodoHandler) Search(c *fiber.Ctx) error {
	userID := mustUserID(c)
	q := c.Query("q")
	if q == "" {
		return fiber.ErrUnprocessableEntity
	}
	todos, err := h.svc.Search(c.Context(), userID, q)
	if err != nil {
		return err
	}
	return c.JSON(todos)
}

func mustUserID(c *fiber.Ctx) uuid.UUID {
	id, _ := uuid.Parse(c.Locals("userID").(string))
	return id
}

func mapErr(err error) error {
	switch err {
	case service.ErrNotFound:
		return fiber.ErrNotFound
	case service.ErrForbidden:
		return fiber.ErrForbidden
	case service.ErrInvalidTransition:
		return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid status transition")
	}
	return err
}
