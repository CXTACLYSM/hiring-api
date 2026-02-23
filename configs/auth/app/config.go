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
	JwtSecret   string

	Http Http
	Grpc Grpc
}

type Http struct {
	Port int

	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
}

type Grpc struct {
	Port int
}

func (c *Config) Validate() error {
	var errorList []error

	if c.Environment != Development && c.Environment != Production {
		errorList = append(errorList, errors.New("environment must be one of: development, production"))
	}
	if c.Name == "" {
		errorList = append(errorList, errors.New("application name can't be empty"))
	}
	if c.Version == "" {
		errorList = append(errorList, errors.New("version can't be empty"))
	}
	if c.Host == "" {
		errorList = append(errorList, errors.New("host can't be empty"))
	}
	if c.Http.Port == 0 {
		errorList = append(errorList, errors.New("http port can't be zero value"))
	}
	if c.Grpc.Port == 0 {
		errorList = append(errorList, errors.New("gRPC port can't be zero value"))
	}
	if c.JwtSecret == "" {
		errorList = append(errorList, errors.New("jwt secret can't be empty"))
	}

	return errors.Join(errorList...)
}

func (c *Config) HttpSocketStr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Http.Port)
}

func (c *Config) GrpcSocketStr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Grpc.Port)
}
