package app

import (
	"errors"
	"fmt"
)

type Config struct {
	Version string
	Host    string
	Port    string
}

func (c *Config) Validate() error {
	var errorList []error

	if c.Version == "" {
		errorList = append(errorList, fmt.Errorf("version can't be empty"))
	}
	if c.Host == "" {
		errorList = append(errorList, fmt.Errorf("host can't be empty"))
	}
	if c.Host == "" {
		errorList = append(errorList, fmt.Errorf("port can't be empty"))
	}

	return errors.Join(errorList...)
}

func (c *Config) Url(protocol string) string {
	return fmt.Sprintf("%s://%s:%s", protocol, c.Host, c.Port)
}

func (c *Config) SocketStr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
