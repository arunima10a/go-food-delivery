package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	customErrors "github.com/arunima10a/go-food-delivery/internal/common/errors"
	"github.com/arunima10a/go-food-delivery/internal/services/catalog-service/internal/products/models"
	"github.com/arunima10a/go-food-delivery/internal/services/catalog-service/internal/products/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type ProductHandler struct {
	repo   repository.ProductRepository
	Logger zerolog.Logger
}

func NewProductHandler(repo repository.ProductRepository, logger zerolog.Logger) *ProductHandler {
	return &ProductHandler{repo: repo, Logger: logger}
}

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description string  `json:"description" validate:"required,max=250"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Category    string  `json:"category" validate:"required"`
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Save a new product to the database
// @Tags products
// @Accept json
// @Produce json
// @Param product body CreateProductRequest true "Product to create"
// @Success 201 {object} models.Product
// @Failure 400 {object} errors.ApiError
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	h.Logger.Info().Msg("Received request to create product")

	req := new(CreateProductRequest)
	if err := c.Bind(req); err != nil {
		return customErrors.NewApiError(http.StatusBadRequest, "Invalid JSON")
	}


	product := &models.Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		CreatedAt:   time.Now(),
	}

	
	event := models.ProductCreatedEvent{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
	}
	eventBytes, _ := json.Marshal(event)

	outbox := &models.OutBoxMessage{
		ID:        uuid.New(),
		Payload:   eventBytes,
		Type:      "ProductCreated",
		Processed: false,
		CreatedAt: time.Now(),
	}

	
	if err := h.repo.Create(product, outbox); err != nil {
		h.Logger.Error().Err(err).Msg("Failed to save product to database")
		return customErrors.NewApiError(http.StatusInternalServerError, "Database Failure")
	}

	h.Logger.Info().Str("productId", product.ID.String()).Msg("Product and Outbox saved successfully")

	return c.JSON(http.StatusCreated, product)
}


func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	products, err := h.repo.GetAll()
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to fetch products")
		return customErrors.NewApiError(http.StatusInternalServerError, "Could not fetch products")
	}
	if products == nil {
		products = []models.Product{}
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetById(c echo.Context) error {
	idStr := c.Param("id")
	productID, err := uuid.Parse(idStr)
	if err != nil {
		return customErrors.NewApiError(http.StatusBadRequest, "Invalid product ID format")
	}
	

	product, err := h.repo.GetById(productID)
	 if err != nil {
		return customErrors.NewApiError(http.StatusNotFound, "Product not found ")
	}

	h.Logger.Info().
		Str("productId", product.ID.String()).
		Msg("GetByID request completed successfully")

	return c.JSON(http.StatusOK, product)

}
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	idStr := c.Param("id")
	productID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product ID"})
	}
	product, err := h.repo.GetById(productID)
	if err != nil {
		return customErrors.NewApiError(http.StatusNotFound, "Product not found")
	}
	

	type UpdateProductRequest struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}
	req := new(UpdateProductRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	

	event := models.ProductCreatedEvent{ 
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
	}
	eventBytes, _ := json.Marshal(event)
	outbox := &models.OutBoxMessage{
		ID:        uuid.New(),
		Payload:   eventBytes,
		Type:      "ProductUpdated",
		Processed: false,
		CreatedAt: time.Now(),
	}

	
	if err := h.repo.UpdateWithOutbox(product, outbox); err != nil {
		return customErrors.NewApiError(http.StatusInternalServerError, "Update failed")
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id, _ := uuid.Parse(c.Param("id"))

	event := map[string]interface{}{"id": id}
	payload, _ := json.Marshal(event)

	outbox := &models.OutBoxMessage{
		ID:        uuid.New(),
		Payload:   payload,
		Type:      "ProductDeleted",
		Processed: false,
		CreatedAt: time.Now(),
	}

	
	if err := h.repo.DeleteWithOutbox(id, outbox); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Delete failed")
	}

	return c.NoContent(http.StatusNoContent)
}
