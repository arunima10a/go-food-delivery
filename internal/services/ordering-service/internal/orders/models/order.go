package models

import (
	"time"
	"github.com/google/uuid"
)

type OrderStaus string

const (
	OrderPending OrderStaus ="PENDING"
	OrderCompleted OrderStaus ="COMPLETED"
	OrderCancelled OrderStaus ="CANCELLED"

)

type Order struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	UserID uuid.UUID `json:"userId"`
	ProductID uuid.UUID `json:"productId"`
	Quantity int `json:"quantity"`
	TotalPrice float64 `json:"totalPrice"`
	Status OrderStaus `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}
