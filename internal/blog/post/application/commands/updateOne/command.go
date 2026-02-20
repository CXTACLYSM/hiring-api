package updateOne

import "github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"

type Command struct {
	Id      string
	Name    string
	Content string
	UserId  string
}

type Handler interface {
	Handle(Command) (*entities.Post, error)
}
