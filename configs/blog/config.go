package configs

import (
	"errors"
	"fmt"
	"time"

	"github.com/CXTACLYSM/hiring-api/configs/blog/app"
	"github.com/CXTACLYSM/hiring-api/configs/blog/database/cache/redis"
	"github.com/CXTACLYSM/hiring-api/configs/blog/database/persistence/postgres"
	"github.com/CXTACLYSM/hiring-api/configs/blog/services"
	"github.com/spf13/viper"
)

type Config struct {
	App             *app.Config
	Auth            *services.Auth
	PostgresCluster *postgres.ClusterConfig
	Redis           *redis.Config
	Kafka           *Kafka
}

type Kafka struct {
	Host string
	Port int
}

func (k *Kafka) Brokers() []string {
	return []string{fmt.Sprintf("%s:%d", k.Host, k.Port)}
}

func (k *Kafka) Validate() error {
	var errorList []error
	if k.Host == "" {
		errorList = append(errorList, fmt.Errorf("kafka host is required"))
	}
	if k.Port == 0 {
		errorList = append(errorList, fmt.Errorf("kafka port is required"))
	}
	return errors.Join(errorList...)
}

func Create() (*Config, error) {
	viper.AutomaticEnv()

	config := &Config{
		App: &app.Config{
			Name:    viper.GetString("APP_NAME"),
			Version: viper.GetString("APP_VERSION"),
			Host:    viper.GetString("APP_HOST"),
			Http: app.Http{
				Port:              viper.GetInt("APP_HTTP_PORT"),
				ReadHeaderTimeout: 5 * time.Second,
				ReadTimeout:       10 * time.Second,
				WriteTimeout:      35 * time.Second,
				IdleTimeout:       60 * time.Second,
				MaxHeaderBytes:    1 << 20,
			},
		},
		Auth: &services.Auth{
			Host:     viper.GetString("AUTH_GRPC_HOST"),
			GrpcPort: viper.GetInt("AUTH_GRPC_PORT"),
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
		Kafka: &Kafka{
			Host: viper.GetString("KAFKA_HOST"),
			Port: viper.GetInt("KAFKA_PORT"),
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
		c.Auth.Validate(),
		c.PostgresCluster.Validate(),
		c.Redis.Validate(),
		c.Kafka.Validate(),
	)
}
