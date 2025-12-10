package configs

import (
	"errors"
	"fmt"

	"github.com/CXTACLYSM/hiring-api/configs/app"
	"github.com/CXTACLYSM/hiring-api/configs/database/analytics/clickhouse"
	"github.com/CXTACLYSM/hiring-api/configs/database/persistence/postgres"
	"github.com/spf13/viper"
)

type Config struct {
	App        app.Config
	Postgres   postgres.Config
	ClickHouse clickhouse.Config
}

func Create() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("configs reading error: %w", err)
		}
	}

	config := &Config{
		App: app.Config{
			Host: viper.GetString("APP_HOST"),
			Port: viper.GetString("APP_PORT"),
		},
		Postgres: postgres.Config{
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetInt("POSTGRES_PORT"),
			Username: viper.GetString("POSTGRES_USERNAME"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			Database: viper.GetString("POSTGRES_DB"),
		},
		ClickHouse: clickhouse.Config{
			Host:     viper.GetString("CLICKHOUSE_HOST"),
			Port:     viper.GetInt("CLICKHOUSE_PORT"),
			Username: viper.GetString("CLICKHOUSE_USERNAME"),
			Password: viper.GetString("CLICKHOUSE_PASSWORD"),
			Database: viper.GetString("CLICKHOUSE_DB"),
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
		c.Postgres.Validate(),
		c.ClickHouse.Validate(),
	)
}
