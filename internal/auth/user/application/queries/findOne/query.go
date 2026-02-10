package findOne

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"
	"github.com/CXTACLYSM/hiring-api/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type Query struct {
	Id       string
	Username string
	Email    string
}

type Handler struct {
	connector *postgres.Connector
}

func NewFindOneUserQueryHandler(connector *postgres.Connector) *Handler {
	return &Handler{
		connector: connector,
	}
}

func (h *Handler) Handle(query Query) (*entities.User, error) {
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

	row := h.connector.ReadPool.QueryRow(context.Background(), sql, args...)

	user := &entities.User{}
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
