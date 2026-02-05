package handlers

import (
	"net/http"
	"time"

	"github.com/arunima10a/go-food-delivery/internal/services/identity-service/config"
	"github.com/arunima10a/go-food-delivery/internal/services/identity-service/internal/users/models"
	"github.com/arunima10a/go-food-delivery/internal/services/identity-service/internal/users/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo repository.UserRepository
	cfg  *config.Config
}

func NewUserHandler(repo repository.UserRepository, cfg *config.Config) *UserHandler {
	return &UserHandler{repo: repo, cfg: cfg}
}

func (h *UserHandler) Register(c echo.Context) error {

	type RegisterRequest struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}
	req := new(RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})

	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Encryption Failed"})
	}
	user := &models.User{
		ID:       uuid.New(),
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}
	if err := h.repo.CreateUser(user); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "User already exists"})

	}
	return c.JSON(http.StatusCreated, user)
}
func (h *UserHandler) Login(c echo.Context) error {

	type LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	user, err := h.repo.FindByEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(h.cfg.JWT.Secret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})

	}
	return c.JSON(http.StatusOK, map[string]string{"token": tokenString})

}
