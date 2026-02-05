package models

import "github.com/google/uuid"

type ProductCreatedEvent struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
}
