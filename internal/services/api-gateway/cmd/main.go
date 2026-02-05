package main

import (
	"log"
	"net/url"

	"github.com/arunima10a/go-food-delivery/internal/services/api-gateway/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.GetConfig()
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	createProxy := func(target string) *middleware.ProxyTarget {
		url, _ := url.Parse(target)
		return &middleware.ProxyTarget{URL: url}
	}

	e.Group("/api/v1/identity").Use(middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		createProxy(cfg.Services.Identity),
	})))

	e.Group("/api/v1/products").Use(middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		createProxy(cfg.Services.Catalog),
	})))

	e.Group("/api/v1/stock").Use(middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		createProxy(cfg.Services.Inventory),
	})))

	e.Group("/api/v1/orders").Use(middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		createProxy(cfg.Services.Ordering),
	})))

	e.Group("/api/v1/search").Use(middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		createProxy(cfg.Services.Search),
	})))

	log.Printf("API Gateway started on %s", cfg.Service.Port)
	e.Logger.Fatal(e.Start(cfg.Service.Port))
}
