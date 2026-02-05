package models

import (
	"time"

	"github.com/google/uuid"
)

type OutBoxMessage struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	Payload   []byte    `gorm:"type:jsonb" json:"paylod"`
	Type      string    `json:"type"`
	Processed bool      `gorm:"index" json"processed"`
	CreatedAt time.Time `json:"createdAt`
}
