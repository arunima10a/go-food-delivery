package main

import (
	"github.com/arunima10a/go-food-delivery/internal/common/logging"
	"github.com/arunima10a/go-food-delivery/internal/services/notification-service/internal/common/messaging"
	"github.com/arunima10a/go-food-delivery/internal/services/notification-service/config"
)

func main() {
	cfg := config.GetConfig()

	logger := logging.NewLogger("notification-service")
	logger.Info().Msg("Notification service is starting...")

	messaging.ConsumeOrderEvents(logger, cfg.RabbitMQ.Host)
}
