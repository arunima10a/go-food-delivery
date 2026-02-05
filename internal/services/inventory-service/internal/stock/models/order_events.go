package models

import "github.com/google/uuid"

type OrderCreatedEvent struct {
	OrderID uuid.UUID `json:"orderId"`
	ProductID uuid.UUID `json:"productId"`
	Quantity int `json:"quantity"`
}

type ProductDeletedEvent struct {
	ID uuid.UUID `json:"id"`
}
