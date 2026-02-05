package database


import(
	"fmt"
	"log"
	"time"
	"github.com/arunima10a/go-food-delivery/internal/services/search-service/config"
	"github.com/arunima10a/go-food-delivery/internal/services/search-service/internal/products/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
func NewPostgresDB(cfg *config.Config)(*gorm.DB, error){
	dsn :=fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
	cfg.Postgres.Host,
	cfg.Postgres.User,
	cfg.Postgres.Password,
	cfg.Postgres.DbName,
	cfg.Postgres.Port,
	cfg.Postgres.SslMode,

	)
	var db *gorm.DB
		var err error
	
		for i := 0; i < 5; i++ {
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err == nil {
				break 
			}
			log.Printf("Database not ready, retrying in 2 seconds... (%d/5)", i+1)
			time.Sleep(2 * time.Second)
		}
	
		if err != nil {
			return nil, err
		}
	
		db.AutoMigrate(&models.ProductSearchModel{}) 
		return db, nil
	}