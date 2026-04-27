package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/todo-app/file-service/internal/service"
)

type FileHandler struct{ svc service.FileService }

func NewFileHandler(svc service.FileService) *FileHandler { return &FileHandler{svc} }

func (h *FileHandler) Upload(c *fiber.Ctx) error {
	userID, _ := uuid.Parse(c.Locals("userID").(string))
	todoID, err := uuid.Parse(c.FormValue("todo_id"))
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}

	file, err := c.FormFile("file")
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}

	f, err := file.Open()
	if err != nil {
		return fiber.ErrInternalServerError
	}
	defer f.Close()

	data := make([]byte, file.Size)
	if _, err := f.Read(data); err != nil {
		return fiber.ErrInternalServerError
	}

	mimeType := file.Header.Get("Content-Type")
	attachment, err := h.svc.Upload(c.Context(), userID, todoID, file.Filename, mimeType, data)
	if err != nil {
		switch err {
		case service.ErrInvalidMIME:
			return fiber.NewError(fiber.StatusUnprocessableEntity, "file type not allowed")
		case service.ErrFileTooLarge:
			return fiber.NewError(fiber.StatusRequestEntityTooLarge, "file exceeds 10MB")
		case service.ErrAttachmentLimit:
			return fiber.NewError(fiber.StatusUnprocessableEntity, "attachment limit reached")
		}
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(attachment)
}

func (h *FileHandler) Download(c *fiber.Ctx) error {
	userID, _ := uuid.Parse(c.Locals("userID").(string))
	fileID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}
	path, err := h.svc.GetPath(c.Context(), userID, fileID)
	if err != nil {
		switch err {
		case service.ErrNotFound:
			return fiber.ErrNotFound
		case service.ErrForbidden:
			return fiber.ErrForbidden
		}
		return err
	}
	return c.SendFile(path)
}

func (h *FileHandler) Delete(c *fiber.Ctx) error {
	userID, _ := uuid.Parse(c.Locals("userID").(string))
	fileID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}
	if err := h.svc.Delete(c.Context(), userID, fileID); err != nil {
		switch err {
		case service.ErrNotFound:
			return fiber.ErrNotFound
		case service.ErrForbidden:
			return fiber.ErrForbidden
		}
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok", "service": "file-service"})
}
