package createOne

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/enums"
	"github.com/CXTACLYSM/hiring-api/pkg/postgres"
	"github.com/google/uuid"
)

type Command struct {
	Username     string
	Email        string
	PasswordHash string
}

type Handler struct {
	connector *postgres.Connector
}

func NewCreateOneUserQueryHandler(connector *postgres.Connector) *Handler {
	return &Handler{
		connector: connector,
	}
}

func (h *Handler) Handle(command Command) (*entities.User, error) {
	now := time.Now()
	user := &entities.User{
		Id:           uuid.NewString(),
		Username:     command.Username,
		Email:        command.Email,
		PasswordHash: command.PasswordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
		DeletedAt:    nil,
	}
	log.Printf("WritePool host: %s\n", h.connector.WritePool.Config().ConnConfig.Host)
	sql := fmt.Sprintf("INSERT INTO %s (id, username, email, password_hash, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", enums.TableUsers)
	_, err := h.connector.WritePool.Exec(
		context.Background(),
		sql,
		user.Id,
		user.Username,
		user.Email,
		command.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
