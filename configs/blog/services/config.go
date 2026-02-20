package services

import (
	"errors"
	"fmt"
)

type Auth struct {
	Host     string
	GrpcPort int
}

func (a *Auth) Validate() error {
	var errorList []error
	if a.Host == "" {
		errorList = append(errorList, fmt.Errorf("auth grpc host is required"))
	}
	if a.GrpcPort == 0 {
		errorList = append(errorList, fmt.Errorf("auth grpc port is required"))
	}
	return errors.Join(errorList...)
}

func (a *Auth) GrpcSocketStr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.GrpcPort)
}
