package deleteOne

import (
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/commands/deleteOne"
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
	return nil
}
