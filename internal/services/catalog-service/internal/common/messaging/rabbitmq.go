package messaging

import (
	"encoding/json"
	"os"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)
func getRabbitMQURL() string {
	host := os.Getenv("RABBITMQ_HOST")
	if host == "" {
		host = "localhost"
	}
	return fmt.Sprintf("amqp://guest:guest@%s:5672/", host)
}

func PublishEvent(exchange string, routingKey string, event interface{}) error {
	conn, err := amqp.Dial(getRabbitMQURL())
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	err = ch.ExchangeDeclare("catalog_events", "fanout", true, false, false, false, nil)
	body, _ := json.Marshal(event)

	return ch.Publish(
		"catalog_events",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
