package messaging

import (
	"github.com/arunima10a/go-food-delivery/internal/services/search-service/internal/products/models"
	"github.com/arunima10a/go-food-delivery/internal/services/search-service/internal/products/repository"

	"encoding/json"
	"log"
	"os"
	"fmt"
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

func ConsumeProductCreated(repo repository.SearchRepository) {
	conn, err := amqp.Dial(getRabbitMQURL())
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	ch.ExchangeDeclare("catalog_events", "fanout", true, false, false, false, nil)

	q, _ := ch.QueueDeclare("search_product_created", false, false, false, false, nil)

	ch.QueueBind(q.Name, "", "catalog_events", false, nil)

	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

	go func() {
		for d := range msgs {
			var deleteEvent models.ProductDeletedEvent
			json.Unmarshal(d.Body, &deleteEvent)

			if deleteEvent.ID != uuid.Nil {
				repo.Delete(deleteEvent.ID)
				log.Printf("[Search Service] Product Details: %s", deleteEvent.ID)
			}
			
			var event models.ProductCreatedEvent
			json.Unmarshal(d.Body, &event)

			searchProduct := &models.ProductSearchModel{
				ID:    event.ID,
				Name:  event.Name,
				Description: event.Description,
				Price: event.Price,
				Category: event.Category,
			}
            if err := repo.Save(searchProduct); err != nil{
				log.Printf("Failed to sync product to search DB: %v", err)
			}else{
			log.Printf("[Search Service] Succesfully Synced: %s", event.Name)
			}
		}
	}()
	select {}
}
