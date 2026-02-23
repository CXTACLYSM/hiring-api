package deleteOne

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/commands/deleteOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/enums"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewDeleteOnePostHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		pool: pool,
	}
}

func (h *Handler) Handle(command deleteOne.Command) error {
	if command.Id == "" {
		return errors.New("id is empty")
	}
	if command.UserId == "" {
		return errors.New("user_id is empty")
	}

	sql := fmt.Sprintf(
		"UPDATE %s SET deleted_at=$1 WHERE id=$2 AND created_by=$3 AND deleted_at IS NULL",
		enums.TablePosts,
	)

	_, err := h.pool.Exec(
		context.Background(),
		sql,
		time.Now(),
		command.Id,
		command.UserId,
	)
	if err != nil {
		return fmt.Errorf("error executing delete post query: %w", err)
	}

	return nil
}
