package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/arunima10a/go-food-delivery/internal/services/inventory-service/internal/stock/models"
	"github.com/arunima10a/go-food-delivery/internal/services/inventory-service/internal/stock/repository"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func getRabbitMQURL() string {
	host := os.Getenv("RABBITMQ_HOST")
	if host == "" {
		host = "localhost" // Fallback for running on Mac
	}
	return fmt.Sprintf("amqp://guest:guest@%s:5672/", host)
}

func ConsumerProductCreated(repo repository.StockRepository) {

	conn, err := amqp.Dial(getRabbitMQURL())
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	ch.ExchangeDeclare("catalog_events", "fanout", true, false, false, false, nil)

	q, err := ch.QueueDeclare(
		"product_created",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
	ch.QueueBind(q.Name, "", "catalog_events", false, nil)

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}
	fmt.Println("!!! INVENTORY CONSUMER IS NOW LISTENING ON catalog_events !!!")

	for d := range msgs {
		fmt.Printf(" INVENTORY RECEIVED A MESSAGE: %s !!!\n", string(d.Body))
		var event models.ProductCreatedEvent
		json.Unmarshal(d.Body, &event)

		stock := models.Stock{
			ID:        uuid.New(),
			ProductID: event.ID,
			Quantity:  10, // Default stock
		}
		repo.CreateStock(&stock)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			var event models.ProductCreatedEvent

			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Printf("Error decoding JSON: %s, err")
				continue
			}
			stock := models.Stock{
				ID:        uuid.New(),
				ProductID: event.ID,
				Quantity:  10,
			}
			if err := repo.CreateStock(&stock); err != nil {
				log.Printf("[Inventory Service] DB Error: %v", err)
			} else {
				log.Printf("[Inventory Service] Success: Stock record created for %s", event.Name)
			}

			log.Printf("[Inventory Service] Received a message")
			log.Printf(" -> Product ID: %s", event.ID)
			log.Printf(" -> Product Name: %s", event.Name)
			log.Printf(" ->Price: $%.2f", event.Price)
			log.Printf(" -------------------------")
		}
	}()
	log.Printf(" [*] Inventory Service waiting for message. To exit press CTRL+C")
	<-forever
}

func ConsumeOrderCreated(repo repository.StockRepository) {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672")
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()

	q, _ := ch.QueueDeclare("order_created", false, false, false, false, nil)
	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

	go func() {
		for d := range msgs {
			var event models.OrderCreatedEvent
			json.Unmarshal(d.Body, &event)

			stock, err := repo.GetStockByProductID(event.ProductID)
			if err == nil {
				stock.Quantity -= event.Quantity

				repo.UpdateStock(stock)
				log.Printf("[Inventory] Stock reduced for product %s. New Total: %d", event.ProductID, stock.Quantity)
			}
		}
	}()
	select {}
}
