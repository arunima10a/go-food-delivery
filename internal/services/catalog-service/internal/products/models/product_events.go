package models

import (
	"github.com/google/uuid"
)

type ProductCreatedEvent struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
}
type ProductDeletedEvent struct {
	ID uuid.UUID `json:"id"`
}
