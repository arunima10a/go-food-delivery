package models

import (
	"time"

	"github.com/google/uuid"
)

type Stock struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	ProductID uuid.UUID `gorm:"uniqueIndex;type:uuid" json:"productId"`
	Quantity  int       `gorm:"default:0" json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
}


