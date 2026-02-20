package services

import (
	"errors"
	"testing"

	"github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_Me_Success(t *testing.T) {
	expectedUser := &entities.User{
		Id:       "test-uuid-123",
		Username: "testuser",
		Email:    "test@email.com",
	}

	finder := &mockFindOneUser{
		user: expectedUser,
		err:  nil,
	}
	service := NewUserService(finder, nil)

	user, err := service.Me("test-uuid-123")

	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserService_Me_NotFound(t *testing.T) {
	finder := &mockFindOneUser{
		user: nil,
		err:  nil,
	}
	service := NewUserService(finder, nil)

	user, err := service.Me("nonexistent-id")

	require.NoError(t, err)
	assert.Nil(t, user)
}

func TestUserService_Me_DatabaseError(t *testing.T) {
	dbErr := errors.New("connection refused")
	finder := &mockFindOneUser{
		user: nil,
		err:  dbErr,
	}
	service := NewUserService(finder, nil)

	user, err := service.Me("test-uuid-123")

	assert.Error(t, err)
	assert.Equal(t, "connection refused", err.Error())
	assert.Nil(t, user)
}
