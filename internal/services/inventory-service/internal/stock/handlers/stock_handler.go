package handlers

import (
	"net/http"
	"github.com/arunima10a/go-food-delivery/internal/services/inventory-service/internal/stock/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type StockHandler struct {
	repo repository.StockRepository
}

func NewStockHandler(repo repository.StockRepository) *StockHandler {
	return &StockHandler{repo: repo}
}
func (h *StockHandler) GetStockByProduct(c echo.Context) error {
	
		productID, err := uuid.Parse(c.Param("productID"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Product ID"})
	}
	stock, err := h.repo.GetStockByProductID(productID)
	if err != nil {
		return  c.JSON(http.StatusNotFound, map[string]string{"error": "Stock not found"})

	}
	return c.JSON(http.StatusOK, stock)
}

func (h *StockHandler) UpdateStock(c echo.Context) error {
	productID, _ := uuid.Parse(c.Param("productID"))

	type UpdateRequest struct {
		Quuatiy int `json:"quantity"`

	}
	req := new(UpdateRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	stock, err := h.repo.GetStockByProductID(productID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Stock not found"})
	}
	stock.Quantity = req.Quuatiy
	if err := h.repo.UpdateStock(stock); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update stock"})
	}
	return c.JSON(http.StatusOK, stock)
}