package createOne

import (
	"context"
	"fmt"

	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/commands/createOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/enums"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewCreateOnePostHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		pool: pool,
	}
}

func (h *Handler) Handle(command createOne.Command) (*entities.Post, error) {
	post := entities.NewPost(command.Name, command.Content, command.UserId)

	sql := fmt.Sprintf("INSERT INTO %s (id, name, content, created_by, updated_by, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", enums.TablePosts)

	_, err := h.pool.Exec(
		context.Background(),
		sql,
		post.Id,
		post.Name,
		post.Content,
		post.CreatedBy,
		post.UpdatedBy,
		post.CreatedAt,
		post.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error executing create one post command: %w", err)
	}

	return post, nil
}
