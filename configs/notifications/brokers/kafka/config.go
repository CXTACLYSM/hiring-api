package kafka

import (
	"errors"
	"fmt"

	"github.com/CXTACLYSM/hiring-api/pkg/events"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (c *Config) Brokers() []string {
	return []string{
		fmt.Sprintf("%s:%d", c.Host, c.Port),
	}
}

func (c *Config) Topics() []string {
	return []string{
		events.TopicUserCreated,
	}
}

func (c *Config) Validate() error {
	var errorList []error
	if c.Host == "" {
		errorList = append(errorList, fmt.Errorf("kafka host is required"))
	}
	if c.Port == 0 {
		errorList = append(errorList, fmt.Errorf("kafka port is required"))
	}

	return errors.Join(errorList...)
}
