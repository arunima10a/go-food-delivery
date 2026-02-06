package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/arunima10a/go-food-delivery/internal/services/ordering-service/config"
	"github.com/arunima10a/go-food-delivery/internal/services/ordering-service/internal/orders/models"
	"github.com/arunima10a/go-food-delivery/internal/services/ordering-service/internal/orders/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	cfg  *config.Config
	repo repository.OrderRepository
	Logger zerolog.Logger
}

func NewOrderHandler(cfg *config.Config, repo repository.OrderRepository, logger zerolog.Logger) *OrderHandler {
	return &OrderHandler{
		cfg:    cfg,
		repo:   repo,
		Logger: logger,
	}
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {

	h.Logger.Info().
		Str("method", c.Request().Method).
		Str("path", c.Path()).
		Msg("New order request received")

	userContext := c.Get("user")
	if userContext == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: No token found"})

	}
	token, ok := userContext.(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error: token type mismatch"})

	}

	claims := token.Claims.(jwt.MapClaims)

	log.Println("DEBUG: Request reached the CreateOrder handler")

	userValue := c.Get("user")
	if userValue == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	userIdStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIdStr)

	type OrderRequest struct {
		ProductID uuid.UUID `json:"productId"`
		Quantity  int       `json:"quantity"`
	}

	req := new(OrderRequest)
	if err := c.Bind(req); err != nil {
		log.Printf("JSON BINDING ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":   "Invalid JSON",
			"details": err.Error(),
		})
	}
	log.Printf("DEBUG: Received request for Product: %s, Qty: %d", req.ProductID, req.Quantity)
	if req.ProductID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ProductID is required"})
	}

	catalogUrl := fmt.Sprintf("%s/api/v1/products/%s", h.cfg.ExternalServices.Catalog, req.ProductID)
	log.Printf("DEBUG: Ordering Service is calling Catalog at: %s", catalogUrl)
	catResp, err := http.Get(catalogUrl)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Catalog service unreachable"})
	}
	defer catResp.Body.Close()

	if catResp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found in Catalog"})
	}

	var product struct {
		Price float64 `json:"price"`
	}
	json.NewDecoder(catResp.Body).Decode(&product)

	inventoryUrl := fmt.Sprintf("%s/api/v1/stock/%s", h.cfg.ExternalServices.Inventory, req.ProductID)
	invResp, err := http.Get(inventoryUrl)

	if err != nil || invResp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Stock info is not available"})
	}
	var stock struct {
		Quantity int `json:"quantity"`
	}
	json.NewDecoder(invResp.Body).Decode(&stock)

	if stock.Quantity < req.Quantity {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Insufficient stock"})
	}

	order := &models.Order{
		ID:         uuid.New(),
		UserID:     userID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		TotalPrice: product.Price * float64(req.Quantity),
		Status:     models.OrderPending,
	}

	fmt.Printf("\n[ORDER SERVICE] >>> Order %s saved successfully. Sending to RabbitMQ...\n", order.ID)

	event := models.OrderCreatedEvent{
		OrderID:   order.ID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}
	eventBytes, _ := json.Marshal(event)

	outbox := &models.OutBoxMessage{
		ID:        uuid.New(),
		Payload:   eventBytes,
		Type:      "OrderCreated",
		Processed: false,
		CreatedAt: time.Now(),
	}

	err = h.repo.CreateOrderWithOutbox(order, outbox)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to place order")
	}
	fmt.Printf("[ORDER SERVICE] Success: Order %s and Outbox record saved to DB.\n", order.ID)

	return c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetMyOrders(c echo.Context) error {

	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println("!!!   GET HISTORY REQUEST RECEIVED NOW      !!!")
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID, _ := uuid.Parse(claims["user_id"].(string))

	orders, err := h.repo.GetOrdersByUserId(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Couls not fetch orders")
	}

	return c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateStatus(c echo.Context) error {
	orderId, _ := uuid.Parse(c.Param("id"))

	type StatusRequest struct {
		Status models.OrderStatus `json:"status`
	}
	req := new(StatusRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	if err := h.repo.UpdateOrderStatus(orderId, req.Status); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failes to update status")
	}
	return c.NoContent(http.StatusOK)
}
