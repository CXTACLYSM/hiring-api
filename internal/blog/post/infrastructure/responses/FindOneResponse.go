package responses

import "github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"

type FindOneResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func NewFindOneResponse(post *entities.Post) *FindOneResponse {
	return &FindOneResponse{
		Id:      post.Id,
		Name:    post.Name,
		Content: post.Content,
	}
}
