package services

import (
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/commands/createOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"
)

type mockFindOneUser struct {
	user *entities.User
	err  error
}

type mockCreateOneUser struct {
	user *entities.User
	err  error
}

func (m *mockFindOneUser) Handle(query findOne.Query) (*entities.User, error) {
	return m.user, m.err
}

func (m *mockCreateOneUser) Handle(command createOne.Command) (*entities.User, error) {
	return m.user, m.err
}
