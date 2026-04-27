package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileAttachment struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TodoID      uuid.UUID `gorm:"type:uuid;not null;index" json:"todo_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Filename    string    `gorm:"not null;size:255" json:"filename"`
	StoragePath string    `gorm:"not null" json:"-"`
	MimeType    string    `gorm:"not null;size:100" json:"mime_type"`
	SizeBytes   int64     `gorm:"not null" json:"size_bytes"`
	CreatedAt   time.Time `json:"created_at"`
}

func (f *FileAttachment) BeforeCreate(_ *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// AllowedMimeTypes defines the allowlist for file uploads (FILE-02).
var AllowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
	"application/pdf": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"text/plain": true,
}

const MaxFileSizeBytes = 10 * 1024 * 1024 // 10MB (FILE-01)
const MaxAttachmentsPerTodo = 10           // FILE-03
