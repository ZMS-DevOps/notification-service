package main

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/mmmajder/zms-devops-notification-service/startup"
	cfg "github.com/mmmajder/zms-devops-notification-service/startup/config"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdin)
	log.SetOutput(os.Stderr)
	log.SetOutput(os.Stdout)
	config := cfg.NewConfig()
	server := startup.NewServer(config)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"security.protocol": "sasl_plaintext",
		"sasl.mechanism":    "PLAIN",
		"sasl.username":     "user1",
		"sasl.password":     config.KafkaAuthPassword,
		"group.id":          "notification-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	consumer.SubscribeTopics([]string{"user.created", "host-review.created", "accommodation-review.created", "reservation-request.created", "reservation.canceled", "host-reviewed-reservation-request"}, nil)
	topicHandlers := map[string]func(*kafka.Message){
		"user.created":                      server.SettingsHandler.OnUserCreated,
		"host-review.created":               server.NotificationHandler.OnHostRated,
		"accommodation-review.created":      server.NotificationHandler.OnAccommodationRated,
		"reservation-request.created":       server.NotificationHandler.OnNewReservationRequestCreated,
		"reservation.canceled":              server.NotificationHandler.OnReservationCancellation,
		"host-reviewed-reservation-request": server.NotificationHandler.OnHostRespondedToReservationRequest,
	}

	go func() {
		for {
			msg, err := consumer.ReadMessage(-1)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}
			if msg == nil {
				log.Printf("Received nil message")
				continue
			}

			handlerFunc, ok := topicHandlers[*msg.TopicPartition.Topic]
			if !ok {
				log.Printf("No handler for topic: %s\n", *msg.TopicPartition.Topic)
				continue
			}
			if handlerFunc == nil {
				log.Printf("Handler function for topic %s is nil", *msg.TopicPartition.Topic)
				continue
			}

			handlerFunc(msg)
		}
	}()

	server.Start()
}
