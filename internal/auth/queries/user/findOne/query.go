package findOne

import (
	"context"
	"errors"

	"github.com/CXTACLYSM/hiring-api/internal/auth/domain/entities"
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
	row := h.connector.ReadPool.QueryRow(
		context.Background(),
		"SELECT id, username, email, password_hash FROM users where email = $1 or username = $2",
		query.Email,
		query.Username,
	)

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
