package messaging

import (
	"encoding/json"
	"os"

	"fmt"

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
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
