package findOne

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewFindOneUserQueryHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		pool: pool,
	}
}

func (h *Handler) Handle(ctx context.Context, query findOne.Query) (*entities.User, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if query.Id != "" {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argIdx))
		args = append(args, query.Id)
		argIdx++
	}
	if query.Email != "" {
		conditions = append(conditions, fmt.Sprintf("email = $%d", argIdx))
		args = append(args, query.Email)
		argIdx++
	}
	if query.Username != "" {
		conditions = append(conditions, fmt.Sprintf("username = $%d", argIdx))
		args = append(args, query.Username)
		argIdx++
	}

	if len(conditions) == 0 {
		return nil, errors.New("at least one search parameter required")
	}

	sql := "SELECT id, username, email, password_hash FROM users WHERE " + strings.Join(conditions, " OR ")

	row := h.pool.QueryRow(ctx, sql, args...)

	user := &entities.User{}
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error scanning row: %w", err)
	}

	return user, nil
}
