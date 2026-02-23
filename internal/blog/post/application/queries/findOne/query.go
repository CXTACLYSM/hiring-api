package findOne

import "github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"

type Query struct {
	UserId string
	Id     string
}

type Handler interface {
	Handle(Query) (*entities.Post, error)
}
