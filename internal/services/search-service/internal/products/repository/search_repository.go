package repository

import (
	"log"

	"github.com/google/uuid"

	"github.com/arunima10a/go-food-delivery/internal/common/utils"
	"github.com/arunima10a/go-food-delivery/internal/services/search-service/internal/products/models"
	"gorm.io/gorm"
)

type SearchRepository interface {
	Save(product *models.ProductSearchModel) error
	AdvancedSearch(name string, category string, minPrice, maxPrice float64, page, pageSize int) (*utils.Pagination, error)
	Delete(id uuid.UUID) error
}
type pgSearchRepository struct {
	db *gorm.DB
}

func NewSearchRepository(db *gorm.DB) SearchRepository {
	return &pgSearchRepository{db: db}
}
func (r *pgSearchRepository) Save(product *models.ProductSearchModel) error {
	return r.db.Save(product).Error
}
func (r *pgSearchRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ProductSearchModel{}, "id = ?", id).Error
}
func (r *pgSearchRepository) AdvancedSearch(name string, category string, minPrice, maxPrice float64, page, pageSize int) (*utils.Pagination, error) {
	var products []models.ProductSearchModel
	var totalItems int64

	query := r.db.Debug().Model(&models.ProductSearchModel{})

	if name != "" {
		searchTerm := "%" + name + "%"
		query = query.Where("name ILIKE ? OR category ILIKE ?", searchTerm, searchTerm)
	}
	if category != "" {
		log.Printf("DEBUG: Applying category filter: %s", category)
		query = query.Where("category ILIKE ?", category)
	}
	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}
	query.Count(&totalItems)

	offset := (page - 1) * pageSize

	err := query.
		Offset(offset).
		Limit(pageSize).
		Order("price ASC").
		Find(&products).Error

	if err != nil {
		return nil, err
	}
	return utils.NewPagination(page, pageSize, totalItems, products), nil
}
