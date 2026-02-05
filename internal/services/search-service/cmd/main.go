package main

import (
	"log"
	"github.com/labstack/echo/v4"
	"github.com/arunima10a/go-food-delivery/internal/services/search-service/config"
	"github.com/arunima10a/go-food-delivery/internal/services/search-service/internal/common/database"
	"github.com/arunima10a/go-food-delivery/internal/services/search-service/internal/products/handlers"
	"github.com/arunima10a/go-food-delivery/internal/services/search-service/internal/common/messaging"
	"github.com/arunima10a/go-food-delivery/internal/services/search-service/internal/products/repository"
)

func main() {
	cfg := config.GetConfig()
	log.Printf("!!! ARCHITECTURE CHECK: Search Service is connecting to: %s on %s:%s", 
        cfg.Postgres.DbName, cfg.Postgres.Host, cfg.Postgres.Port)


	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	searchRepo := repository.NewSearchRepository(db)
	SearchHandler := handlers.NewSearchHandler(searchRepo)

	go func() {
		messaging.ConsumeProductCreated(searchRepo)
	}()
	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "Search Service is Alive")
	})

	e.GET("/api/v1/search", SearchHandler.Search)

	log.Printf("Search Service running on %s", cfg.Service.Port)
	e.Logger.Fatal(e.Start(cfg.Service.Port))


}
