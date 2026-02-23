package configs

import (
	"errors"
	"time"

	"github.com/CXTACLYSM/hiring-api/configs/auth/app"
	"github.com/CXTACLYSM/hiring-api/configs/auth/brokers/kafka"
	"github.com/CXTACLYSM/hiring-api/configs/auth/database/cache/redis"
	"github.com/CXTACLYSM/hiring-api/configs/auth/database/persistence/postgres"
	"github.com/spf13/viper"
)

type Config struct {
	App             *app.Config
	PostgresCluster *postgres.ClusterConfig
	Redis           *redis.Config
	Kafka           *kafka.Config
}

func Create() (*Config, error) {
	viper.AutomaticEnv()

	config := &Config{
		App: &app.Config{
			Environment: viper.GetString("APP_ENVIRONMENT"),
			Name:        viper.GetString("APP_NAME"),
			Version:     viper.GetString("APP_VERSION"),
			Host:        viper.GetString("APP_HOST"),
			JwtSecret:   viper.GetString("JWT_SECRET"),

			Http: app.Http{
				Port:              viper.GetInt("APP_HTTP_PORT"),
				ReadHeaderTimeout: 5 * time.Second,
				ReadTimeout:       10 * time.Second,
				WriteTimeout:      35 * time.Second,
				IdleTimeout:       60 * time.Second,
				MaxHeaderBytes:    1 << 20,
			},
			Grpc: app.Grpc{
				Port: viper.GetInt("APP_GRPC_PORT"),
			},
		},
		PostgresCluster: &postgres.ClusterConfig{
			Read: &postgres.Config{
				Host:     viper.GetString("POSTGRES_READ_HOST"),
				Port:     viper.GetInt("POSTGRES_READ_PORT"),
				Username: viper.GetString("POSTGRES_READ_USERNAME"),
				Password: viper.GetString("POSTGRES_READ_PASSWORD"),
				Database: viper.GetString("POSTGRES_READ_DATABASE"),
			},
			Write: &postgres.Config{
				Host:     viper.GetString("POSTGRES_WRITE_HOST"),
				Port:     viper.GetInt("POSTGRES_WRITE_PORT"),
				Username: viper.GetString("POSTGRES_WRITE_USERNAME"),
				Password: viper.GetString("POSTGRES_WRITE_PASSWORD"),
				Database: viper.GetString("POSTGRES_WRITE_DATABASE"),
			},
		},
		Redis: &redis.Config{
			Host:             viper.GetString("REDIS_HOST"),
			Port:             viper.GetInt("REDIS_PORT"),
			Username:         viper.GetString("REDIS_USERNAME"),
			Password:         viper.GetString("REDIS_PASSWORD"),
			AuthDatabase:     viper.GetInt("REDIS_AUTH_DB"),
			ResourceDatabase: viper.GetInt("REDIS_RESOURCE_DB"),
		},
		Kafka: &kafka.Config{
			Host:     viper.GetString("KAFKA_HOST"),
			Port:     viper.GetInt("KAFKA_PORT"),
			Username: viper.GetString("KAFKA_USERNAME"),
			Password: viper.GetString("KAFKA_PASSWORD"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	return errors.Join(
		c.App.Validate(),
		c.PostgresCluster.Validate(),
		c.Redis.Validate(),
		c.Kafka.Validate(),
	)
}
