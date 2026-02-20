package responses

import "github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"

func NewFindListResponse(postList []*entities.Post) []*FindOneResponse {
	list := make([]*FindOneResponse, 0, len(postList))

	for _, post := range postList {
		list = append(list, NewFindOneResponse(post))
	}

	return list
}
