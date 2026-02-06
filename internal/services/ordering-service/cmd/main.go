package main

import (
	"log"

	"github.com/arunima10a/go-food-delivery/internal/services/ordering-service/config"
	"github.com/arunima10a/go-food-delivery/internal/services/ordering-service/internal/common/database"
	"github.com/arunima10a/go-food-delivery/internal/services/ordering-service/internal/orders/handlers"
	"github.com/arunima10a/go-food-delivery/internal/services/ordering-service/internal/orders/repository"
	"github.com/arunima10a/go-food-delivery/internal/services/ordering-service/internal/orders/worker"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.GetConfig()

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal(err)

	}

	orderRepo := repository.NewPostgresOrderRepository(db)
	orderHandler := handlers.NewOrderHandler(cfg, orderRepo)

	worker.StartOutboxWorker(db)

	e := echo.New()

	jwtConfig := echojwt.Config{
		SigningKey: []byte(cfg.JWT.Secret),
	}

	r := e.Group("/api/v1/orders")
	r.Use(echojwt.WithConfig(jwtConfig))

	r.POST("", orderHandler.CreateOrder)
	r.GET("", orderHandler.GetMyOrders)
	r.PUT("/:id/status", orderHandler.UpdateStatus)

	log.Printf("Ordering service started on port %s", cfg.Service.Port)
	e.Logger.Fatal(e.Start(cfg.Service.Port))
}
