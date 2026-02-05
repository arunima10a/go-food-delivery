package main

import (
	"fmt"
	"net/http"

	"github.com/arunima10a/go-food-delivery/internal/common/logging"
	"github.com/arunima10a/go-food-delivery/internal/common/middleware"
	_ "github.com/arunima10a/go-food-delivery/internal/services/catalog-service/docs"

	customMiddleware "github.com/arunima10a/go-food-delivery/internal/common/middleware"
	"github.com/arunima10a/go-food-delivery/internal/services/catalog-service/config"
	"github.com/arunima10a/go-food-delivery/internal/services/catalog-service/internal/common/database"

	CustomValidator "github.com/arunima10a/go-food-delivery/internal/services/catalog-service/internal/common/validator"
	"github.com/arunima10a/go-food-delivery/internal/services/catalog-service/internal/products/handlers"
	"github.com/arunima10a/go-food-delivery/internal/services/catalog-service/internal/products/models"
	"github.com/arunima10a/go-food-delivery/internal/services/catalog-service/internal/products/repository"
	"github.com/arunima10a/go-food-delivery/internal/services/catalog-service/internal/products/worker"
	"github.com/go-playground/validator/v10"
	echojwt "github.com/labstack/echo-jwt/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
)

// @title Food Delivery catalog API
// @version 1.0
// @description This is the Catalog Service for our food delivery App.
// @host localhost:5001
// @BasePath /api/v1

func main() {
	cfg := config.GetConfig()
	logger := logging.NewLogger("catalog-service")

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not connect to database")
	}

	err = db.AutoMigrate(&models.Product{},&models.OutBoxMessage{})
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not migrate database")
	}
	logger.Info().Msg("Database connection successful and migrated")

	productRepo := repository.NewProductRepository(db)

	worker.StartOutboxWorker(db, logger)

	productHandler := handlers.NewProductHandler(productRepo, logger)

	e := echo.New()

	e.HTTPErrorHandler = customMiddleware.CustomHTTPErrorHandler(logger)

	jwtConfig := echojwt.Config{
		SigningKey: []byte(cfg.JWT.Secret),
	}

	adminGroup := e.Group("/api/v1/products")
	adminGroup.Use(echojwt.WithConfig(jwtConfig))
	adminGroup.Use(middleware.RequireRole("admin"))

	adminGroup.POST("", productHandler.CreateProduct)
	adminGroup.DELETE("/:id", productHandler.DeleteProduct)
	adminGroup.PUT("/:id", productHandler.UpdateProduct)

	e.GET("/api/v1/products", productHandler.GetAllProducts)
	e.GET("/api/v1/products/:id", productHandler.GetById)
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("%s is Running", cfg.Service.Name))
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Validator = &CustomValidator.CustomValidator{Validator: validator.New()}

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("%s is Running", cfg.Service.Name))

	})

	logger.Info().
		Str("port", cfg.Service.Port).
		Str("evn", "development").
		Msgf("Catalog service is starting")
	e.Logger.Fatal(e.Start(cfg.Service.Port))

}
