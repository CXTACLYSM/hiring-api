package services

import (
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"
)

type TokenGenerator interface {
	Generate(user *entities.User) (string, error)
}
