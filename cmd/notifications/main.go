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
	"go.uber.org/zap"
)

func main() {
	cfg, err := configs.Create()
	if err != nil {
		log.Fatalf("error creating config: %v", err)
	}

	container := &di.Container{}
	if err = container.Init(cfg); err != nil {
		log.Fatalf("error initializing container: %v", err)
	}

	logger := container.Infrastructure.Logger
	defer logger.Sync()
	defer container.Infrastructure.PgConnector.Close()
	defer container.Infrastructure.RedisConnector.Close()

	defer func() {
		if err := container.Infrastructure.Kafka.ConsumerGroup.Close(); err != nil {
			logger.Error("error closing consumer group", zap.Error(err))
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for err := range container.Infrastructure.Kafka.ConsumerGroup.Errors() {
			logger.Error("consumer group error", zap.Error(err))
		}
	}()

	go func() {
		for {
			logger.Debug("attempting to join consumer group...")
			if err := container.Infrastructure.Kafka.ConsumerGroup.Consume(ctx, container.Infrastructure.Kafka.Topics, container.Handlers.UserCreatedHandler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					logger.Info("consumer group closed")
					return
				}
				logger.Error("consumer group consume error", zap.Error(err))
			}

			if ctx.Err() != nil {
				logger.Info("context cancelled, stopping consumer")
				return
			}

			logger.Info("rebalance happened, rejoining consumer group...")
		}
	}()

	logger.Info("notifications service started",
		zap.Strings("topics", container.Infrastructure.Kafka.Topics),
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down notifications service...")
	cancel()
	logger.Info("notifications service stopped")
}
