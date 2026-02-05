package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arunima10a/go-food-delivery/internal/services/catalog-service/internal/products/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetById(id uuid.UUID) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)

}

func (m *MockProductRepository) Create(p *models.Product, o *models.OutBoxMessage) error {
	args := m.Called(p, o)
	return args.Error(0)
}

func (m *MockProductRepository) GetAll() ([]models.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) DeleteWithOutbox(id uuid.UUID, o *models.OutBoxMessage) error {
	args := m.Called(id, o)
	return args.Error(0)
}

func (m *MockProductRepository) UpdateWithOutbox(p *models.Product, o *models.OutBoxMessage) error {
	args := m.Called(p, o)
	return args.Error(0)
}

func TestGetProductByID_Success(t *testing.T) {
	e := echo.New()
	productID := uuid.New()
	mockRepo := new(MockProductRepository)

	expectedProduct := &models.Product{
		ID:    productID,
		Name:  "Test Burger",
		Price: 10.0,
	}

	mockRepo.On("GetById", productID).Return(expectedProduct, nil)

	h := NewProductHandler(mockRepo, zerolog.Nop())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/api/v1/products/:id")
	c.SetParamNames("id")
	c.SetParamValues(productID.String())

	err := h.GetById(c)
	
	assert.NoError(t, err) 
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), productID.String());
	
}
