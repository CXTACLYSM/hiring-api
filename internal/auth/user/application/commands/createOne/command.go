package createOne

import "github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"

type Command struct {
	Username     string
	Email        string
	PasswordHash string
}

type Handler interface {
	Handle(Command) (*entities.User, error)
}
