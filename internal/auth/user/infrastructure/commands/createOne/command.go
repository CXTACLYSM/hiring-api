package createOne

import (
	"context"
	"fmt"

	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/commands/createOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/infrastructure/enums"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewCreateOneUserQueryHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		pool: pool,
	}
}

func (h *Handler) Handle(command createOne.Command) (*entities.User, error) {
	user := entities.NewUser(command.Username, command.Email, command.PasswordHash)

	sql := fmt.Sprintf("INSERT INTO %s (id, username, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", enums.TableUsers)
	_, err := h.pool.Exec(
		context.Background(),
		sql,
		user.Id,
		user.Username,
		user.Email,
		command.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error handling create user query: %w", err)
	}

	return user, nil
}
