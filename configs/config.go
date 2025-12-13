package configs

import (
	"errors"

	"github.com/CXTACLYSM/hiring-api/configs/app"
	"github.com/CXTACLYSM/hiring-api/configs/database/persistence/postgres"
	"github.com/spf13/viper"
)

type Config struct {
	App             *app.Config
	PostgresCluster *postgres.ClusterConfig
}

func Create() (*Config, error) {
	viper.AutomaticEnv()

	config := &Config{
		App: &app.Config{
			Version: viper.GetString("APP_VERSION"),
			Host:    viper.GetString("APP_HOST"),
			Port:    viper.GetString("APP_PORT"),
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
	)
}
