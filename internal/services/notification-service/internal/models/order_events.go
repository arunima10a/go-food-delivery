package models

import "github.com/google/uuid"

type OrderCreatedEvent struct {
	OrderID    uuid.UUID `json:"orderId"`
	ProductID  uuid.UUID `json:"productId"`
	UserID     uuid.UUID `json:"userId"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"totalPrice"`
}
