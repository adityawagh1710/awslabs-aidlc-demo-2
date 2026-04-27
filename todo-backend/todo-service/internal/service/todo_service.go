package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
	"github.com/todo-app/todo-service/internal/model"
	"github.com/todo-app/todo-service/internal/repository"
)

var (
	ErrNotFound          = errors.New("todo not found")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidTransition = errors.New("invalid status transition")
)

type TodoService interface {
	Create(ctx context.Context, userID uuid.UUID, input CreateTodoInput) (*model.Todo, error)
	List(ctx context.Context, userID uuid.UUID, f repository.TodoFilter) ([]model.Todo, error)
	Get(ctx context.Context, userID, todoID uuid.UUID) (*model.Todo, error)
	Update(ctx context.Context, userID, todoID uuid.UUID, input UpdateTodoInput) (*model.Todo, error)
	Delete(ctx context.Context, userID, todoID uuid.UUID) error
	Search(ctx context.Context, userID uuid.UUID, query string) ([]model.Todo, error)
}

type CreateTodoInput struct {
	Title       string
	Description string
	Priority    model.TodoPriority
	DueDate     *time.Time
	TagIDs      []uuid.UUID
}

type UpdateTodoInput struct {
	Title       *string
	Description *string
	Status      *model.TodoStatus
	Priority    *model.TodoPriority
	DueDate     *time.Time
	TagIDs      []uuid.UUID
}

type todoService struct {
	todoRepo    repository.TodoRepository
	tagRepo     repository.TagRepository
	outboxRepo  repository.OutboxRepository
	es          *elasticsearch.Client
	fileSvcURL  string
	schedSvcURL string
}

func NewTodoService(
	todoRepo repository.TodoRepository,
	tagRepo repository.TagRepository,
	outboxRepo repository.OutboxRepository,
	es *elasticsearch.Client,
	fileSvcURL, schedSvcURL string,
) TodoService {
	return &todoService{todoRepo, tagRepo, outboxRepo, es, fileSvcURL, schedSvcURL}
}

func (s *todoService) Create(ctx context.Context, userID uuid.UUID, input CreateTodoInput) (*model.Todo, error) {
	tags, err := s.resolveTags(ctx, userID, input.TagIDs)
	if err != nil {
		return nil, err
	}
	todo := &model.Todo{
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
		Status:      model.StatusPending,
		Tags:        tags,
	}
	if err := s.todoRepo.Insert(ctx, todo); err != nil {
		return nil, err
	}
	s.enqueueOutbox(ctx, todo, "upsert")
	return todo, nil
}

func (s *todoService) List(ctx context.Context, userID uuid.UUID, f repository.TodoFilter) ([]model.Todo, error) {
	return s.todoRepo.FindByUser(ctx, userID, f)
}

func (s *todoService) Get(ctx context.Context, userID, todoID uuid.UUID) (*model.Todo, error) {
	todo, err := s.todoRepo.FindByID(ctx, todoID)
	if err != nil {
		return nil, ErrNotFound
	}
	if todo.UserID != userID {
		return nil, ErrForbidden
	}
	return todo, nil
}

func (s *todoService) Update(ctx context.Context, userID, todoID uuid.UUID, input UpdateTodoInput) (*model.Todo, error) {
	todo, err := s.Get(ctx, userID, todoID)
	if err != nil {
		return nil, err
	}

	if input.Status != nil {
		if err := validateTransition(todo.Status, *input.Status); err != nil {
			return nil, err
		}
		todo.Status = *input.Status
	}
	if input.Title != nil {
		todo.Title = *input.Title
	}
	if input.Description != nil {
		todo.Description = *input.Description
	}
	if input.Priority != nil {
		todo.Priority = *input.Priority
	}
	if input.DueDate != nil {
		todo.DueDate = input.DueDate
	}
	if input.TagIDs != nil {
		tags, err := s.resolveTags(ctx, userID, input.TagIDs)
		if err != nil {
			return nil, err
		}
		todo.Tags = tags
	}

	if err := s.todoRepo.Update(ctx, todo); err != nil {
		return nil, err
	}
	s.enqueueOutbox(ctx, todo, "upsert")
	return todo, nil
}

func (s *todoService) Delete(ctx context.Context, userID, todoID uuid.UUID) error {
	todo, err := s.Get(ctx, userID, todoID)
	if err != nil {
		return err
	}
	if err := s.todoRepo.SoftDelete(ctx, todo.ID); err != nil {
		return err
	}
	s.enqueueOutbox(ctx, todo, "delete")
	return nil
}

func (s *todoService) Search(ctx context.Context, userID uuid.UUID, query string) ([]model.Todo, error) {
	// Build ES query scoped to userID, excluding deleted
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":  query,
						"fields": []string{"title^2", "description"},
					},
				},
				"filter": []map[string]interface{}{
					{"term": map[string]interface{}{"user_id": userID.String()}},
					{"bool": map[string]interface{}{"must_not": map[string]interface{}{"exists": map[string]interface{}{"field": "deleted_at"}}}},
				},
			},
		},
	}
	body, _ := json.Marshal(esQuery)
	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex("todos"),
		s.es.Search.WithBody(jsonReader(body)),
	)
	if err != nil {
		// Fallback: return empty (caller can degrade gracefully)
		return nil, nil
	}
	defer res.Body.Close()
	return parseESHits(res)
}

// ValidateTransition enforces TODO-04: pending→in_progress→done only.
func ValidateTransition(current, next model.TodoStatus) error {
	return validateTransition(current, next)
}

func validateTransition(current, next model.TodoStatus) error {
	if current == next {
		return nil // no-op, same status
	}
	allowed, ok := model.ValidTransitions[current]
	if !ok {
		return ErrInvalidTransition
	}
	for _, s := range allowed {
		if s == next {
			return nil
		}
	}
	return ErrInvalidTransition
}

func (s *todoService) resolveTags(ctx context.Context, userID uuid.UUID, tagIDs []uuid.UUID) ([]model.Tag, error) {
	var tags []model.Tag
	for _, id := range tagIDs {
		tag, err := s.tagRepo.FindByID(ctx, id)
		if err != nil || tag.UserID != userID {
			return nil, ErrForbidden
		}
		tags = append(tags, *tag)
	}
	return tags, nil
}

func (s *todoService) enqueueOutbox(ctx context.Context, todo *model.Todo, op string) {
	payload, _ := json.Marshal(todo)
	_ = s.outboxRepo.Insert(ctx, &model.SearchOutbox{
		TodoID:    todo.ID,
		Operation: op,
		Payload:   payload,
	})
}
