package worker

import (
	"encoding/json"
	"log"
	"time"

	"github.com/arunima10a/go-food-delivery/internal/services/ordering-service/internal/common/messaging"
	"github.com/arunima10a/go-food-delivery/internal/services/ordering-service/internal/orders/models"
	"gorm.io/gorm"
)

func StartOutboxWorker(db *gorm.DB) {
	go func() {
		for {
			var messages []models.OutBoxMessage

			db.Where("processed = ?", false).Find(&messages)

			for _, msg := range messages {
				var event models.OrderCreatedEvent
				json.Unmarshal(msg.Payload, &event)

				err := messaging.PublishEvent("order_events", "", event)
				if err == nil {
					db.Model(&msg).Update("processed", true)
					log.Printf("[Outbox Worker] Successfully sent message %s", msg.ID)
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()
}
