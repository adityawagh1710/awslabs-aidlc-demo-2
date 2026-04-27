package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/todo-app/file-service/internal/model"
	"github.com/todo-app/file-service/internal/repository"
)

var (
	ErrNotFound        = errors.New("file not found")
	ErrForbidden       = errors.New("forbidden")
	ErrInvalidMIME     = errors.New("file type not allowed")
	ErrFileTooLarge    = errors.New("file exceeds 10MB limit")
	ErrAttachmentLimit = errors.New("attachment limit reached")
)

type FileService interface {
	Upload(ctx context.Context, userID, todoID uuid.UUID, filename, mimeType string, data []byte) (*model.FileAttachment, error)
	GetPath(ctx context.Context, userID, fileID uuid.UUID) (string, error)
	Delete(ctx context.Context, userID, fileID uuid.UUID) error
}

type fileService struct {
	repo        repository.FileRepository
	storagePath string
}

func NewFileService(repo repository.FileRepository, storagePath string) FileService {
	return &fileService{repo, storagePath}
}

func (s *fileService) Upload(ctx context.Context, userID, todoID uuid.UUID, filename, mimeType string, data []byte) (*model.FileAttachment, error) {
	if !model.AllowedMimeTypes[mimeType] {
		return nil, ErrInvalidMIME
	}
	if int64(len(data)) > model.MaxFileSizeBytes {
		return nil, ErrFileTooLarge
	}
	count, err := s.repo.CountByTodo(ctx, todoID)
	if err != nil {
		return nil, err
	}
	if count >= model.MaxAttachmentsPerTodo {
		return nil, ErrAttachmentLimit
	}

	// FILE-06: uploads/{userID}/{todoID}/{uuid}_{filename}
	storagePath := filepath.Join(s.storagePath, userID.String(), todoID.String(),
		fmt.Sprintf("%s_%s", uuid.New().String(), filepath.Base(filename)))

	if err := os.MkdirAll(filepath.Dir(storagePath), 0750); err != nil {
		return nil, err
	}
	if err := os.WriteFile(storagePath, data, 0640); err != nil {
		return nil, err
	}

	f := &model.FileAttachment{
		TodoID:      todoID,
		UserID:      userID,
		Filename:    filepath.Base(filename),
		StoragePath: storagePath,
		MimeType:    mimeType,
		SizeBytes:   int64(len(data)),
	}
	return f, s.repo.Insert(ctx, f)
}

func (s *fileService) GetPath(ctx context.Context, userID, fileID uuid.UUID) (string, error) {
	f, err := s.repo.FindByID(ctx, fileID)
	if err != nil {
		return "", ErrNotFound
	}
	if f.UserID != userID {
		return "", ErrForbidden
	}
	return f.StoragePath, nil
}

func (s *fileService) Delete(ctx context.Context, userID, fileID uuid.UUID) error {
	f, err := s.repo.FindByID(ctx, fileID)
	if err != nil {
		return ErrNotFound
	}
	if f.UserID != userID {
		return ErrForbidden
	}
	_ = os.Remove(f.StoragePath)
	return s.repo.Delete(ctx, fileID)
}
