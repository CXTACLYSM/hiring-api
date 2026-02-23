package updateOne

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/commands/updateOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/enums"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewUpdateOnePostHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		pool: pool,
	}
}

func (h *Handler) Handle(command updateOne.Command) (*entities.Post, error) {
	var setClauses []string
	var whereClauses []string
	var args []any
	argIdx := 1

	if command.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name=$%d", argIdx))
		args = append(args, command.Name)
		argIdx++
	}
	if command.Content != "" {
		setClauses = append(setClauses, fmt.Sprintf("content=$%d", argIdx))
		args = append(args, command.Content)
		argIdx++
	}

	if len(setClauses) == 0 {
		return nil, errors.New("nothing to update")
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at=$%d", argIdx))
	args = append(args, time.Now())
	argIdx++

	whereClauses = append(whereClauses, fmt.Sprintf("id=$%d", argIdx))
	args = append(args, command.Id)
	argIdx++

	whereClauses = append(whereClauses, fmt.Sprintf("created_by=$%d", argIdx))
	args = append(args, command.UserId)
	argIdx++

	sql := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s RETURNING id, name, content, created_at, updated_at, created_by, updated_by",
		enums.TablePosts,
		strings.Join(setClauses, ", "),
		strings.Join(whereClauses, " AND "),
	)

	post := &entities.Post{}
	err := h.pool.QueryRow(context.Background(), sql, args...).Scan(
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
		return nil, fmt.Errorf("error updating post: %w", err)
	}

	return post, nil
}
