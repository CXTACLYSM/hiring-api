package services

import (
	"context"
	"fmt"

	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/commands/createOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"
	"go.uber.org/zap"
)

type UserService struct {
	findOneUser   findOne.Handler
	createOneUser createOne.Handler
	logger        *zap.Logger
}

func NewUserService(findOne findOne.Handler, createOne createOne.Handler, logger *zap.Logger) *UserService {
	return &UserService{
		findOneUser:   findOne,
		createOneUser: createOne,
		logger:        logger.Named("user_service"),
	}
}

func (s *UserService) Me(ctx context.Context, userId string) (*entities.User, error) {
	user, err := s.findOneUser.Handle(ctx, findOne.Query{
		Id: userId,
	})
	if err != nil {
		s.logger.Error("error finding one user by id",
			zap.String("user_id", userId),
			zap.Error(err),
		)
		return nil, fmt.Errorf("error handling find one user query: %w", err)
	}

	return user, nil
}
