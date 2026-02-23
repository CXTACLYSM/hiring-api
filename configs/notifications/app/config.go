package app

import (
	"errors"
	"fmt"
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
}

func (c *Config) Validate() error {
	var errorList []error

	if c.Environment != Development && c.Environment != Production {
		errorList = append(errorList, fmt.Errorf("environment must be one of: %s, %s", Development, Production))
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

	return errors.Join(errorList...)
}
