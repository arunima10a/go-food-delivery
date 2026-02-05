package models

import (
	"time"

	"github.com/google/uuid"
)

type OutBoxMessage struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Payload   []byte    `gorm:"type:jsonb"`
	Type      string
	Processed bool `gorm:"index"`
	CreatedAt time.Time
}
