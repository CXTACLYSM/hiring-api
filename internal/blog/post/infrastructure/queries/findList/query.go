package findList

import (
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/queries/findList"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"
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
	return nil, nil
}
