package findOne

import (
	"context"

	"github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"
)

type Query struct {
	Id       string
	Username string
	Email    string
}

type Handler interface {
	Handle(context.Context, Query) (*entities.User, error)
}
