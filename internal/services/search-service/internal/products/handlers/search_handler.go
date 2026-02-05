package handlers

import (
	"net/http"
	"strconv"

	"github.com/arunima10a/go-food-delivery/internal/services/search-service/internal/products/repository"
	"github.com/labstack/echo/v4"
)

type SearchHandler struct {
	repo repository.SearchRepository

}

func NewSearchHandler(repo repository.SearchRepository) *SearchHandler{ 
	return &SearchHandler{repo: repo}
} 

func (h *SearchHandler) Search(c echo.Context) error {
	name := c.QueryParam("q")
	category := c.QueryParam("category")

	minPrice, _ := strconv.ParseFloat(c.QueryParam("minPrice"), 64)
	maxPrice, _ := strconv.ParseFloat(c.QueryParam("maxPrice"), 64)

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {page = 1}

	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageSize <=0 { pageSize = 10}

	pagination, err := h.repo.AdvancedSearch(name, category, minPrice, maxPrice, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error":"Search failed"})
	} 
	return c.JSON(http.StatusOK, pagination)

	products, err := h.repo.AdvancedSearch(name, category, minPrice, maxPrice, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error":"Search failed"})

	}
	return c.JSON(http.StatusOK, products)
}
