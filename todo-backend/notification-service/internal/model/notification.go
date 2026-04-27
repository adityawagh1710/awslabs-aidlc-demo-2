package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	TodoID    *uuid.UUID `gorm:"type:uuid" json:"todo_id,omitempty"`
	Message   string    `gorm:"not null;size:500" json:"message"`
	Delivered bool      `gorm:"default:false" json:"delivered"`
	CreatedAt time.Time `json:"created_at"`
}

func (n *Notification) BeforeCreate(_ *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
