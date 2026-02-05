package worker

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/arunima10a/go-food-delivery/internal/services/catalog-service/internal/common/messaging"
	"github.com/arunima10a/go-food-delivery/internal/services/catalog-service/internal/products/models"
)

func StartOutboxWorker(db *gorm.DB, logger zerolog.Logger){
	go func() {
		for {
			var messages []models.OutBoxMessage
			db.Where("processed = ?", false).Limit(10).Find(&messages)

			for _, msg := range messages {

				var event interface{}
				json.Unmarshal(msg.Payload, &event)

				err := messaging.PublishEvent("catalog_events", "", msg.Payload)
				if err == nil {
					db.Model(&msg).Update("processed", true)
					logger.Info().Msgf("Succesfully relayed catalog event: %s", msg.ID)
				}else{
					logger.Error().Err(err).Msg("Failed to relay outbox event")
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()
}