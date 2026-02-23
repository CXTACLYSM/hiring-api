package app

import (
	"errors"
	"fmt"
	"time"
)

const (
	Development = "development"
	Production  = "production"
)

type Config struct {
	Environment string
	Name        string
	Version     string
	Host        string
	Http        Http
}

type Http struct {
	Port              int
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
}

func (c *Config) Validate() error {
	var errorList []error

	if c.Environment != Development && c.Environment != Production {
		errorList = append(errorList, errors.New("environment must be one of: development, production"))
	}
	if c.Name == "" {
		errorList = append(errorList, fmt.Errorf("application name can't be empty"))
	}
	if c.Version == "" {
		errorList = append(errorList, fmt.Errorf("version can't be empty"))
	}
	if c.Host == "" {
		errorList = append(errorList, fmt.Errorf("host can't be empty"))
	}
	if c.Http.Port == 0 {
		errorList = append(errorList, fmt.Errorf("http port can't be zero value"))
	}

	return errors.Join(errorList...)
}

func (c *Config) HttpSocketStr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Http.Port)
}
