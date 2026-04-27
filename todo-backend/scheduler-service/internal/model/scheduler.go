package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reminder struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TodoID    uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	FireAt    time.Time `gorm:"not null;index"`
	Fired     bool      `gorm:"default:false"`
	CreatedAt time.Time
}

func (r *Reminder) BeforeCreate(_ *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type RecurrenceConfig struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	TodoID         uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	CronExpression string    `gorm:"not null"`
	NextOccurrence time.Time `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (r *RecurrenceConfig) BeforeCreate(_ *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
