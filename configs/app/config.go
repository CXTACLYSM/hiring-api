package app

import (
	"errors"
	"fmt"
)

type Config struct {
	Host string
	Port string
}

func (c *Config) Validate() error {
	var errorList []error

	if c.Host == "" {
		errorList = append(errorList, fmt.Errorf("url can't be empty"))
	}
	if c.Host == "" {
		errorList = append(errorList, fmt.Errorf("port can't be empty"))
	}

	return errors.Join(errorList...)
}

func (c *Config) Url(protocol string) string {
	return fmt.Sprintf("%s://%s:%s", protocol, c.Host, c.Port)
}
