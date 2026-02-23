package responses

import "github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"

type CreateOneResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func NewCreateOneResponse(post *entities.Post) *CreateOneResponse {
	return &CreateOneResponse{
		Id:      post.Id,
		Name:    post.Name,
		Content: post.Content,
	}
}
