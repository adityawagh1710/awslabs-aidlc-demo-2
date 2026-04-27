package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/todo-app/todo-service/internal/model"
	"github.com/todo-app/todo-service/internal/repository"
)

var ErrTagNotFound = errors.New("tag not found")

type TagService interface {
	Create(ctx context.Context, userID uuid.UUID, name string) (*model.Tag, error)
	List(ctx context.Context, userID uuid.UUID) ([]model.Tag, error)
	Delete(ctx context.Context, userID, tagID uuid.UUID) error
}

type tagService struct{ repo repository.TagRepository }

func NewTagService(repo repository.TagRepository) TagService { return &tagService{repo} }

func (s *tagService) Create(ctx context.Context, userID uuid.UUID, name string) (*model.Tag, error) {
	tag := &model.Tag{UserID: userID, Name: name}
	return tag, s.repo.Insert(ctx, tag)
}

func (s *tagService) List(ctx context.Context, userID uuid.UUID) ([]model.Tag, error) {
	return s.repo.FindByUser(ctx, userID)
}

func (s *tagService) Delete(ctx context.Context, userID, tagID uuid.UUID) error {
	tag, err := s.repo.FindByID(ctx, tagID)
	if err != nil {
		return ErrTagNotFound
	}
	if tag.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(ctx, tagID)
}
