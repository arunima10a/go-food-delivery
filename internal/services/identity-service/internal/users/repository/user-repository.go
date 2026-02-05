package repository

import (
	"github.com/arunima10a/go-food-delivery/internal/services/identity-service/internal/users/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
}

type pgUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) UserRepository {
	return &pgUserRepository{db: db}
}
func (r *pgUserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}
func (r *pgUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}
