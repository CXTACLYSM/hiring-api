package di

import (
	"fmt"
	"log"
	"os"

	"github.com/CXTACLYSM/hiring-api/configs/notifications"
	"github.com/CXTACLYSM/hiring-api/internal/notifications/consumer"
	pgConnector "github.com/CXTACLYSM/hiring-api/pkg/postgres"
	redisConnector "github.com/CXTACLYSM/hiring-api/pkg/redis"
	"github.com/IBM/sarama"
)

type Infrastructure struct {
	PgConnector    *pgConnector.Connector
	RedisConnector *redisConnector.Connector
	Kafka          *Kafka
}

type Kafka struct {
	ConsumerGroup sarama.ConsumerGroup
	Topics        []string
}

type Handlers struct {
	UserCreatedHandler sarama.ConsumerGroupHandler
}

type Container struct {
	Infrastructure *Infrastructure
	Handlers       *Handlers
}

func (c *Container) Init(cfg *configs.Config) error {
	if err := c.initInfrastructure(cfg); err != nil {
		return err
	}
	if err := c.initHandlers(); err != nil {
		return err
	}

	return nil
}

func (c *Container) initInfrastructure(cfg *configs.Config) error {
	readDSN, err := cfg.PostgresCluster.DSN(pgConnector.ReadOperation)
	if err != nil {
		return fmt.Errorf("error initializing infra read pgx pool: %w", err)
	}
	writeDSN, err := cfg.PostgresCluster.DSN(pgConnector.WriteOperation)
	if err != nil {
		return fmt.Errorf("error initializing infra write pgx pool: %w", err)
	}
	pgConn, err := pgConnector.NewConnector(readDSN, writeDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	authCfg, resourceCfg := cfg.Redis.ConnectorConfigs()
	redisConn, err := redisConnector.NewConnector(authCfg, resourceCfg)
	if err != nil {
		return fmt.Errorf("failed to connecto to redis: %w", err)
	}

	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	config := sarama.NewConfig()
	config.Version = sarama.DefaultVersion
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	log.Printf("app name: %s", cfg.App.Name)
	log.Printf("brokers: %v", cfg.Kafka.Brokers())
	group, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers(), cfg.App.Name, config)
	if err != nil {
		return fmt.Errorf("error creating consumer group: %w", err)
	}

	c.Infrastructure = &Infrastructure{
		PgConnector:    pgConn,
		RedisConnector: redisConn,
		Kafka: &Kafka{
			ConsumerGroup: group,
			Topics:        cfg.Kafka.Topics(),
		},
	}

	return nil
}

func (c *Container) initHandlers() error {
	c.Handlers = &Handlers{
		UserCreatedHandler: consumer.NewUserCreatedHandler(),
	}

	return nil
}
