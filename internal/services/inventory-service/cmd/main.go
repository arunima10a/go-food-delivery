package main

import (
	"log"
	"net/http"

	"github.com/arunima10a/go-food-delivery/internal/services/inventory-service/config"
	"github.com/arunima10a/go-food-delivery/internal/services/inventory-service/internal/common/database"
	"github.com/arunima10a/go-food-delivery/internal/services/inventory-service/internal/stock/handlers"
	"github.com/arunima10a/go-food-delivery/internal/services/inventory-service/internal/stock/repository"
	"github.com/labstack/echo/v4"
	"github.com/arunima10a/go-food-delivery/internal/services/inventory-service/internal/common/messaging"
)

func main() {
	cfg := config.GetConfig()

	db, _ := database.NewInventoryDB(cfg)
	
	stockRepo := repository.NewPoatgresStockRepository(db)

	stockHandler := handlers.NewStockHandler(stockRepo)
	

	go func() {
		messaging.ConsumerProductCreated(stockRepo)
	}()
	e := echo.New()

	e.GET("/api/v1/stock/:productID", stockHandler.GetStockByProduct)
	e.PUT("/api/v1/stock/:productID", stockHandler.UpdateStock)
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Inventory Service is Running")
	})
	log.Printf("Inventory HTTP server Starting on port %s", cfg.Service.Port)
	e.Logger.Fatal(e.Start(cfg.Service.Port))

	go func() {
		messaging.ConsumerProductCreated(stockRepo)
	}()

	go func() {
		messaging.ConsumeOrderCreated(stockRepo)
	}()
}
