package repository

import (
	"github.com/arunima10a/go-food-delivery/internal/services/inventory-service/internal/stock/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockRepository interface {
	GetStockByProductID(productID uuid.UUID) (*models.Stock, error)
	UpdateStock(stock *models.Stock) error
	CreateStock(stock *models.Stock) error
}
type pgStockRepository struct {
	db *gorm.DB
}

func NewPoatgresStockRepository(db *gorm.DB) StockRepository {
	return &pgStockRepository{db: db}
}

func (r *pgStockRepository) GetStockByProductID(productID uuid.UUID) (*models.Stock, error) {
	var stock models.Stock
	if err := r.db.First(&stock, "product_id = ?", productID).Error; err != nil {
		return nil, err
	}
	return &stock, nil
}
func (r *pgStockRepository) UpdateStock(stock *models.Stock) error {
	return r.db.Save(stock).Error
}

func (r *pgStockRepository) CreateStock(stock *models.Stock) error {
	return r.db.Create(stock).Error
}