package findList

import (
	"context"
	"fmt"
	"strings"

	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/queries/findList"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/enums"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewFindListPostHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		pool: pool,
	}
}

func (h *Handler) Handle(query findList.Query) ([]*entities.Post, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if query.UserId == "" {
		return nil, fmt.Errorf("user_id cannot be empty")
	}
	if query.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIdx))
		args = append(args, "%"+query.Name+"%")
		argIdx++
	}
	if query.Content != "" {
		conditions = append(conditions, fmt.Sprintf("content ILIKE $%d", argIdx))
		args = append(args, "%"+query.Content+"%")
		argIdx++
	}

	conditions = append(conditions, fmt.Sprintf("created_by=$%d", argIdx))
	args = append(args, query.UserId)
	argIdx++

	sql := fmt.Sprintf(
		"SELECT id, name, content, created_at FROM %s WHERE %s AND deleted_at IS NULL",
		enums.TablePosts,
		strings.Join(conditions, " AND "),
	)
	rows, err := h.pool.Query(
		context.Background(),
		sql,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("error executing find post list query: %w", err)
	}
	defer rows.Close()

	result := make([]*entities.Post, 0)
	for rows.Next() {
		post := &entities.Post{}
		err := rows.Scan(
			&post.Id,
			&post.Name,
			&post.Content,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error mapping posts row: err=%w", err)
		}
		result = append(result, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return result, nil
}
