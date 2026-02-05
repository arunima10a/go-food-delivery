package messaging

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/arunima10a/go-food-delivery/internal/services/notification-service/internal/models"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

func getRabbitMQURL() string {
	host := os.Getenv("RABBITMQ_HOST")
	if host == "" {
		host = "localhost"
	}
	return fmt.Sprintf("amqp://guest:guest@%s:5672/", host)

}

func ConsumeOrderEvents(logger zerolog.Logger, rabbitHost string) {

	dialAddr := fmt.Sprintf("amqp://guest:guest@%s:5672/", rabbitHost)

	logger.Info().Str("host", dialAddr).Msg("Attempting to connect to RabbitMQ...")

	conn, err := amqp.Dial(dialAddr)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to RabbitMQ")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to open a channel")
	}
	defer ch.Close()

	ch.ExchangeDeclare("order_events", "fanout", true, false, false, false, nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("FAILED to declare exchange")
	}

	q, err := ch.QueueDeclare("order_created", false, false, false, false, nil)
	if err != nil {

		logger.Fatal().Err(err).Msg("Failed to declare a queue")
	}

	ch.QueueBind(q.Name, "", "order_events", false, nil)
	if err != nil{
		logger.Fatal().Err(err).Msg("FAILED tp declare queue")
	}



	logger.Info().Msg(" [*] Waiting for order messages. To exit press CTRL+C")

	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

	for d := range msgs {
		var event models.OrderCreatedEvent
		if err := json.Unmarshal(d.Body, &event); err != nil {
			logger.Error().Err(err).Msg("Error decoding JSON")
			continue
		}

		logger.Info().
			Str("orderId", event.OrderID.String()).
			Str("userId", event.UserID.String()).
			Float64("amount", event.TotalPrice).
			Msg(" NOTIFICATION: Sending Order Confirmation Email to User...")

		fmt.Printf("\n--- EMAIL SENT ---\nTo: User %s\nSubject: Order %s Confirmed\nTotal: $%.2f\n----------------\n\n", event.UserID, event.OrderID, event.TotalPrice)
	}

}
