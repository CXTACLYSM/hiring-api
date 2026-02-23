package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	configs "github.com/CXTACLYSM/hiring-api/configs/notifications"
	"github.com/CXTACLYSM/hiring-api/internal/notifications/di"
	"github.com/IBM/sarama"
)

func main() {
	cfg, err := configs.Create()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}

	container := &di.Container{}
	if err = container.Init(cfg); err != nil {
		log.Fatalf("Error initializing container: %s", err.Error())
	}
	defer container.Infrastructure.RedisConnector.Close()
	defer container.Infrastructure.PgConnector.Close()

	defer func() {
		if err := container.Infrastructure.Kafka.ConsumerGroup.Close(); err != nil {
			log.Printf("error closing consumer group: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for err := range container.Infrastructure.Kafka.ConsumerGroup.Errors() {
			log.Printf("consumer group error: %v", err)
		}
	}()

	go func() {
		for {
			log.Println("attempting to join consumer group...")
			if err := container.Infrastructure.Kafka.ConsumerGroup.Consume(ctx, container.Infrastructure.Kafka.Topics, container.Handlers.UserCreatedHandler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					log.Println("consumer group closed")
					return
				}
				log.Printf("consumer group consume error: %v", err)
			}

			if ctx.Err() != nil {
				log.Println("context cancelled, stopping consumer")
				return
			}

			log.Println("rebalance happened, rejoining consumer group...")
		}
	}()

	log.Printf("notifications service started, consuming topics: %v", container.Infrastructure.Kafka.Topics)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down notifications service...")
	cancel()
	log.Println("notifications service stopped")
}
