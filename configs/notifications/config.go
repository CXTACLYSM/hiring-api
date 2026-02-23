package configs

import (
	"errors"

	"github.com/CXTACLYSM/hiring-api/configs/notifications/app"
	"github.com/CXTACLYSM/hiring-api/configs/notifications/brokers/kafka"
	"github.com/CXTACLYSM/hiring-api/configs/notifications/database/cache/redis"
	"github.com/CXTACLYSM/hiring-api/configs/notifications/database/persistence/postgres"
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
			Name:    viper.GetString("APP_NAME"),
			Version: viper.GetString("APP_VERSION"),
			Host:    viper.GetString("APP_HOST"),
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
