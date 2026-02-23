package app

import (
	"errors"
	"fmt"
)

type Config struct {
	Name    string
	Version string
	Host    string
}

func (c *Config) Validate() error {
	var errorList []error

	if c.Name == "" {
		errorList = append(errorList, fmt.Errorf("application name can't be empty"))
	}
	if c.Version == "" {
		errorList = append(errorList, fmt.Errorf("version can't be empty"))
	}
	if c.Host == "" {
		errorList = append(errorList, fmt.Errorf("host can't be empty"))
	}

	return errors.Join(errorList...)
}
