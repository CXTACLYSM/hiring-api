package entities

import (
	"time"

	entities "github.com/CXTACLYSM/hiring-api/pkg/shared/domain"
	"github.com/google/uuid"
)

type User entities.User

func NewUser(username, email, passwordHash string) *User {
	now := time.Now()
	return &User{
		Id:           uuid.NewString(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (u *User) ToShared() *entities.User {
	return (*entities.User)(u)
}
