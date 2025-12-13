package postgres

import (
	"errors"
	"fmt"
)

const (
	WriteOperation = 1
	ReadOperation  = 2
)

type ClusterConfig struct {
	Read  *Config
	Write *Config
}

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func (c *ClusterConfig) DSN(operation uint8) (string, error) {
	switch operation {
	case WriteOperation:
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			c.Read.Host, c.Read.Port, c.Read.Username, c.Read.Password, c.Read.Database,
		), nil
	case ReadOperation:
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			c.Write.Host, c.Write.Port, c.Write.Username, c.Write.Password, c.Write.Database,
		), nil
	default:
		return "", errors.New("invalid pool type")
	}
}

func (c *ClusterConfig) Validate() error {
	var errorList []error

	if c.Read.Host == "" {
		errorList = append(errorList, fmt.Errorf("read host is required"))
	}
	if c.Read.Port == 0 {
		errorList = append(errorList, fmt.Errorf("read port is required"))
	}
	if c.Read.Username == "" {
		errorList = append(errorList, fmt.Errorf("read username is required"))
	}
	if c.Read.Database == "" {
		errorList = append(errorList, fmt.Errorf("read database is required"))
	}

	if c.Write.Host == "" {
		errorList = append(errorList, fmt.Errorf("write host is required"))
	}
	if c.Write.Port == 0 {
		errorList = append(errorList, fmt.Errorf("write port is required"))
	}
	if c.Write.Username == "" {
		errorList = append(errorList, fmt.Errorf("write username is required"))
	}
	if c.Write.Database == "" {
		errorList = append(errorList, fmt.Errorf("write database is required"))
	}

	return errors.Join(errorList...)
}
