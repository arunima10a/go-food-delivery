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
	err := r.db.Where("user_id = ?", userID).Find(&orders).Error
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
