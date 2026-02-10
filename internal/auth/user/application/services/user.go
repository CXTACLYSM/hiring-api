package services

import (
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/commands/createOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"
)

type UserService struct {
	findOneUser   *findOne.Handler
	createOneUser *createOne.Handler
}

func NewUserService(findOne *findOne.Handler, createOne *createOne.Handler) *UserService {
	return &UserService{
		findOneUser:   findOne,
		createOneUser: createOne,
	}
}

func (s *UserService) Me(userId string) (*entities.User, error) {
	user, err := s.findOneUser.Handle(
		findOne.Query{
			Id: userId,
		},
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
