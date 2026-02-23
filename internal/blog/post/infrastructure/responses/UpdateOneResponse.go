package responses

import "github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"

type UpdateOneResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func NewUpdateOneResponse(post *entities.Post) *UpdateOneResponse {
	return &UpdateOneResponse{
		Id:      post.Id,
		Name:    post.Name,
		Content: post.Content,
	}
}
