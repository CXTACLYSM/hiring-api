package updateOne

import "github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"

type Command struct {
	UserId  string
	Id      string
	Name    string
	Content string
}

type Handler interface {
	Handle(Command) (*entities.Post, error)
}
