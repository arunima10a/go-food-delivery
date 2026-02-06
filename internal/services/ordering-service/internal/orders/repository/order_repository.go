package repository

import (
	"github.com/arunima10a/go-food-delivery/internal/services/ordering-service/internal/orders/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	SaveOrder(order *models.Order) error
	GetOrdersByUserId(userID uuid.UUID) ([]models.Order, error)
	CreateOrderWithOutbox(order *models.Order, outbox *models.OutBoxMessage) error
	UpdateOrderStatus(orderId uuid.UUID, status models.OrderStatus) error
}

type pgOrderRepository struct {
	db *gorm.DB
}

func NewPostgresOrderRepository(db *gorm.DB) OrderRepository {
	return &pgOrderRepository{db: db}
}

func (r *pgOrderRepository) SaveOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *pgOrderRepository) GetOrdersByUserId(userID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("user_id = ? AND status != ?", userID, models.OrderArchived).Find(&orders).Error
	return orders, err
}

func (r *pgOrderRepository) CreateOrderWithOutbox(order *models.Order, outbox *models.OutBoxMessage) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		if err := tx.Create(outbox).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *pgOrderRepository) UpdateOrderStatus(orderId uuid.UUID, status models.OrderStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderId).Update("status", status).Error
}
