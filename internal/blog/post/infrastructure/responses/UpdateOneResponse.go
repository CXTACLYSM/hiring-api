package responses

import "github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"

type UpdateOneResponse struct {
	Id      string
	Name    string
	Content string
}

func NewUpdateOneResponse(post *entities.Post) *UpdateOneResponse {
	return &UpdateOneResponse{
		Id:      post.Id,
		Name:    post.Name,
		Content: post.Content,
	}
}
