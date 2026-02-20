package findOne

import "github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"

type Query struct {
	Id       string
	Username string
	Email    string
}

type Handler interface {
	Handle(Query) (*entities.User, error)
}
