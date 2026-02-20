package responses

import "github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"

type CreateOneResponse struct {
	Id      string
	Name    string
	Content string
}

func NewCreateOneResponse(post *entities.Post) *CreateOneResponse {
	return &CreateOneResponse{
		Id:      post.Id,
		Name:    post.Name,
		Content: post.Content,
	}
}
