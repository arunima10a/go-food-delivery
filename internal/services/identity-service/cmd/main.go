package main

import (
	"log"

	"github.com/arunima10a/go-food-delivery/internal/services/identity-service/config"
	"github.com/arunima10a/go-food-delivery/internal/services/identity-service/internal/common/database"
	"github.com/arunima10a/go-food-delivery/internal/services/identity-service/internal/users/handlers"
	"github.com/arunima10a/go-food-delivery/internal/services/identity-service/internal/users/repository"
	"github.com/labstack/echo/v4"
	
)

func main() {
	cfg := config.GetConfig()

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()


	userRepo := repository.NewPostgresUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo, cfg)

	e.POST("/api/v1/identity/register", userHandler.Register)
	e.POST("/api/v1/identity/login", userHandler.Login)
	log.Printf("Identity service running on port %s", cfg.Service.Port)

	e.Logger.Fatal(e.Start(cfg.Service.Port))

}
