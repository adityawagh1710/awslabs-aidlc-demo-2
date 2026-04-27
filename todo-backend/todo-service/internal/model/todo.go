package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoStatus string
type TodoPriority string

const (
	StatusPending    TodoStatus = "pending"
	StatusInProgress TodoStatus = "in_progress"
	StatusDone       TodoStatus = "done"

	PriorityLow    TodoPriority = "low"
	PriorityMedium TodoPriority = "medium"
	PriorityHigh   TodoPriority = "high"
)

// ValidTransitions defines allowed status transitions (TODO-04).
// Forward: pending → in_progress → done. Reopen: done → pending.
var ValidTransitions = map[TodoStatus][]TodoStatus{
	StatusPending:    {StatusInProgress},
	StatusInProgress: {StatusDone},
	StatusDone:       {StatusPending},
}

type Todo struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Title       string         `gorm:"not null;size:255" json:"title"`
	Description string         `gorm:"size:5000" json:"description"`
	Status      TodoStatus     `gorm:"not null;default:pending" json:"status"`
	Priority    TodoPriority   `gorm:"not null;default:medium" json:"priority"`
	DueDate     *time.Time     `json:"due_date,omitempty"`
	Tags        []Tag          `gorm:"many2many:todo_tags;" json:"tags,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (t *Todo) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type Tag struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string    `gorm:"not null;size:50" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *Tag) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// SearchOutbox stores pending Elasticsearch sync events (outbox pattern).
type SearchOutbox struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TodoID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Operation string    `gorm:"not null"` // "upsert" | "delete"
	Payload   []byte    `gorm:"type:jsonb"`
	Processed bool      `gorm:"default:false;index"`
	CreatedAt time.Time
}

// Migration creates the singular table; override GORM's auto-pluralization.
func (SearchOutbox) TableName() string { return "search_outbox" }

func (s *SearchOutbox) BeforeCreate(_ *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
