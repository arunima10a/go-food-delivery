package models

import "github.com/google/uuid"

type ProductSearchModel struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	Name        string    `gorm:"index" json:"name"`
	Description string    `json:"description"`
	Price       float64   `gorm:"index" json:"price"`
	Category    string    `gorm:"index" json:"category"`
}

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
