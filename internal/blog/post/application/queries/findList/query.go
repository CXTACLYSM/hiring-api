package findList

import "github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"

type Query struct {
	UserId  string
	Name    string
	Content string
}

type Handler interface {
	Handle(Query) ([]*entities.Post, error)
}
