package redis

import (
	"errors"
	"fmt"
	"time"

	pkgRedis "github.com/CXTACLYSM/hiring-api/pkg/redis"
)

type Config struct {
	Host             string
	Port             int
	Username         string
	Password         string
	AuthDatabase     int
	ResourceDatabase int
}

func (c *Config) ConnectorConfigs() (authCfg, resourceCfg *pkgRedis.Config) {
	base := pkgRedis.Config{
		Host:           c.Host,
		Port:           c.Port,
		User:           c.Username,
		Password:       c.Password,
		ConnectTimeout: 5 * time.Second,
		ReadTimeout:    3 * time.Second,
		WriteTimeout:   3 * time.Second,
	}

	auth := base
	auth.Db = c.AuthDatabase

	resource := base
	resource.Db = c.ResourceDatabase

	return &auth, &resource
}

func (c *Config) Validate() error {
	var errorList []error

	if c.Host == "" {
		errorList = append(errorList, fmt.Errorf("redis host is required"))
	}
	if c.Port == 0 {
		errorList = append(errorList, fmt.Errorf("redis port is required"))
	}
	if c.Password == "" {
		errorList = append(errorList, fmt.Errorf("redis password is required"))
	}

	return errors.Join(errorList...)
}
