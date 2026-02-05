package repository

import (
	"github.com/arunima10a/go-food-delivery/internal/services/catalog-service/internal/products/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product, outbox *models.OutBoxMessage) error
	GetById(id uuid.UUID) (*models.Product, error)
	GetAll() ([]models.Product, error)
	UpdateWithOutbox(product *models.Product, outbox *models.OutBoxMessage) error
	DeleteWithOutbox(id uuid.UUID, outbox *models.OutBoxMessage) error
}

type pgProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &pgProductRepository{db: db}
}

func (r *pgProductRepository) Create(product *models.Product, outbox *models.OutBoxMessage) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			return err
		}
		if err := tx.Create(outbox).Error; err != nil {
			return err
		}
		return nil
	})
}
func (r *pgProductRepository) GetById(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := r.db.First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *pgProductRepository) GetAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error
	return products, err
}
func (r *pgProductRepository) UpdateWithOutbox(product *models.Product, outbox *models.OutBoxMessage) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(product).Error; err != nil {
			return err
		}
		if err := tx.Create(outbox).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *pgProductRepository) DeleteWithOutbox(id uuid.UUID, outbox *models.OutBoxMessage) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Product{}, "id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Create(outbox).Error; err != nil {
			return err
		}
		return nil
	})
}
