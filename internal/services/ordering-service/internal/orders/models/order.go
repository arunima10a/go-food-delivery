package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderPending   OrderStatus = "PENDING"
	OrderCompleted OrderStatus = "COMPLETED"
	OrderCancelled OrderStatus = "CANCELLED"
	OrderArchived OrderStatus = "ARCHIVED"
)

type Order struct {
	ID         uuid.UUID   `gorm:"primaryKey;type:uuid" json:"id"`
	UserID     uuid.UUID   `json:"userId"`
	ProductID  uuid.UUID   `json:"productId"`
	Quantity   int         `json:"quantity"`
	TotalPrice float64     `json:"totalPrice"`
	Status     OrderStatus `json:"status"`
	CreatedAt  time.Time   `json:"createdAt"`
}
