package findOne

import (
	"context"
	"errors"
	"fmt"

	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/enums"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewFindOnePostHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		pool: pool,
	}
}

func (h *Handler) Handle(query findOne.Query) (*entities.Post, error) {
	sql := fmt.Sprintf(
		"SELECT id, name, content, created_at, updated_at, created_by, updated_by FROM %s WHERE id = $1 AND created_by=$2 AND deleted_at IS NULL",
		enums.TablePosts,
	)

	row := h.pool.QueryRow(
		context.Background(),
		sql,
		query.Id,
		query.UserId,
	)

	post := &entities.Post{}
	err := row.Scan(
		&post.Id,
		&post.Name,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.CreatedBy,
		&post.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error executing find one post query: %w", err)
	}

	return post, nil
}
