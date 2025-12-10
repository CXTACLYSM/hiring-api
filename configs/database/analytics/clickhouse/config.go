package clickhouse

import (
	"errors"
	"fmt"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"tcp://%s:%d?username=%s&password=%s&database=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database,
	)
}

func (c *Config) Validate() error {
	var errorList []error

	if c.Host == "" {
		errorList = append(errorList, fmt.Errorf("host is required"))
	}
	if c.Port == 0 {
		errorList = append(errorList, fmt.Errorf("port is required"))
	}
	if c.Username == "" {
		errorList = append(errorList, fmt.Errorf("username is required"))
	}
	if c.Database == "" {
		errorList = append(errorList, fmt.Errorf("database is required"))
	}

	return errors.Join(errorList...)
}
