package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/todo-app/todo-service/internal/service"
)

type TagHandler struct{ svc service.TagService }

func NewTagHandler(svc service.TagService) *TagHandler { return &TagHandler{svc} }

type createTagRequest struct {
	Name string `json:"name" validate:"required,max=50"`
}

func (h *TagHandler) Create(c *fiber.Ctx) error {
	var req createTagRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	if err := validate.Struct(req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	userID := mustUserID(c)
	tag, err := h.svc.Create(c.Context(), userID, req.Name)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(tag)
}

func (h *TagHandler) List(c *fiber.Ctx) error {
	tags, err := h.svc.List(c.Context(), mustUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(tags)
}

func (h *TagHandler) Delete(c *fiber.Ctx) error {
	tagID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}
	if err := h.svc.Delete(c.Context(), mustUserID(c), tagID); err != nil {
		return mapErr(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
