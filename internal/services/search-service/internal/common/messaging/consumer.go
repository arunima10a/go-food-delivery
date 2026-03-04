package messaging

import (
	"encoding/base64"
	"encoding/json"

	"github.com/arunima10a/go-food-delivery/internal/services/search-service/internal/products/models"
	"github.com/arunima10a/go-food-delivery/internal/services/search-service/internal/products/ai"
	"github.com/arunima10a/go-food-delivery/internal/services/search-service/internal/products/repository"

	"fmt"
	"log"
	"os"

	"github.com/google/uuid"

	"github.com/streadway/amqp"
)

func getRabbitMQURL() string {
	host := os.Getenv("RABBITMQ_HOST")
	if host == "" {
		host = "localhost"
	}
	return fmt.Sprintf("amqp://guest:guest@%s:5672/", host)
}

func ConsumeProductCreated(repo repository.SearchRepository, aiClient *ai.OpenRouterClient) {
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

			var encodedString string
			json.Unmarshal(d.Body, &encodedString)

			decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
			if err != nil {
				log.Printf("FATAL: Base64 Decode Failed: %v", err)
				continue
			}

			var deleteEvent models.ProductDeletedEvent
			json.Unmarshal(d.Body, &deleteEvent)

			if deleteEvent.ID != uuid.Nil {
				repo.Delete(deleteEvent.ID)
				log.Printf("[Search Service] Product Details: %s", deleteEvent.ID)
			}

			var event models.ProductCreatedEvent
			if err := json.Unmarshal(decodedBytes, &event); err != nil {
				log.Printf("FATAL: Final JSON Unmarshal Failed: %v", err)
				continue
			}

			if event.ID == uuid.Nil {
				log.Printf("ERROR: Received message with empty ID. Check your JSON tags!")
				continue
			}

			log.Printf("DEBUG: Consumer received product: %s with description: %s", event.Name, event.Description)

			tags, err := aiClient.GenerateSemanticTags(event.Name, event.Description)
			if err != nil {
				log.Printf("AI Enrichment failed: %v", err)
				tags = "error-fetching-tags"
			}

			log.Printf("DEBUG: AI generated these tags: %s", tags)

			searchProduct := &models.ProductSearchModel{
				ID:          event.ID,
				Name:        event.Name,
				Description: event.Description,
				Price:       event.Price,
				Category:    event.Category,
				AiMetadata:  tags,
			}
			if err := repo.Save(searchProduct); err != nil {
				log.Printf("Failed to sync product to search DB: %v", err)
			} else {
				log.Printf("[Search Service] Succesfully Synced: %s", event.Name)
			}
		}
	}()
	select {}
}
