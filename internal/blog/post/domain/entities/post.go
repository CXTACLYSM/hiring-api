package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`

	CreatedBy string  `json:"created_by"`
	UpdatedBy string  `json:"updated_by"`
	DeletedBy *string `json:"deleted_by"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func NewPost(name, content, userId string) *Post {
	now := time.Now()
	return &Post{
		Id:        uuid.NewString(),
		Name:      name,
		Content:   content,
		CreatedBy: userId,
		UpdatedBy: userId,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
