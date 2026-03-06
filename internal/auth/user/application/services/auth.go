package services

import (
	"context"
	"fmt"

	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/commands/createOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/dto"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/pkg/events"
	"github.com/CXTACLYSM/hiring-api/pkg/kafka"
	appErrors "github.com/CXTACLYSM/hiring-api/pkg/shared/application/errors"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/validation"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	validator      *validation.Validator
	findOneUser    findOne.Handler
	createOneUser  createOne.Handler
	tokenGenerator TokenGenerator
	kafkaPublisher kafka.EventPublisher
	logger         *zap.Logger
}

func NewAuthService(
	validator *validation.Validator,
	findOneUser findOne.Handler,
	createOneUser createOne.Handler,
	tokenGenerator TokenGenerator,
	kafkaPublisher kafka.EventPublisher,
	logger *zap.Logger,
) *AuthService {
	return &AuthService{
		validator:      validator,
		findOneUser:    findOneUser,
		createOneUser:  createOneUser,
		tokenGenerator: tokenGenerator,
		kafkaPublisher: kafkaPublisher,
		logger:         logger.Named("auth_service"),
	}
}

func (s *AuthService) Register(ctx context.Context, dto *dto.RegisterDTO) (string, error) {
	if err := s.validator.Struct(dto); err != nil {
		s.logger.Debug("validation failed",
			zap.String("username", dto.Username),
			zap.String("email", dto.Email),
			zap.Error(err),
		)
		return "", fmt.Errorf("error validating register dto: %w", err)
	}

	user, err := s.findOneUser.Handle(ctx, findOne.Query{
		Username: dto.Username,
		Email:    dto.Email,
	})
	if err != nil {
		s.logger.Error("failed to find user",
			zap.String("username", dto.Username),
			zap.String("email", dto.Email),
			zap.Error(err),
		)
		return "", fmt.Errorf("error handling find one user query: %w", err)
	}
	if user != nil {
		s.logger.Warn("registration attempt for existing user",
			zap.String("username", dto.Username),
			zap.String("email", dto.Email),
		)
		return "", &appErrors.ApplicationError{
			Message: "user already exists",
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", zap.Error(err))
		return "", fmt.Errorf("error generating password hash: %w", err)
	}

	user, err = s.createOneUser.Handle(createOne.Command{
		Username:     dto.Username,
		Email:        dto.Email,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		s.logger.Error("failed to create user",
			zap.String("username", dto.Username),
			zap.Error(err),
		)
		return "", fmt.Errorf("error handling create one user query: %w", err)
	}

	token, err := s.tokenGenerator.Generate(user)
	if err != nil {
		s.logger.Error("failed to generate token",
			zap.String("user_id", user.Id),
			zap.Error(err),
		)
		return "", fmt.Errorf("error generating token: %w", err)
	}

	err = s.kafkaPublisher.Push(&events.UserCreated{
		UserId:    user.Id,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
	if err != nil {
		s.logger.Warn("failed to publish user.created event",
			zap.String("user_id", user.Id),
			zap.String("topic", events.TopicUserCreated),
			zap.Error(err),
		)
	} else {
		s.logger.Info("user registered",
			zap.String("user_id", user.Id),
			zap.String("username", user.Username),
		)
	}

	return token, nil
}

func (s *AuthService) Login(ctx context.Context, dto *dto.LoginDTO) (string, error) {
	if err := s.validator.Struct(dto); err != nil {
		s.logger.Debug("login validation failed",
			zap.String("login", dto.Login),
			zap.Error(err),
		)
		return "", fmt.Errorf("error validation login dto: %w", err)
	}

	user, err := s.findOneUser.Handle(ctx, findOne.Query{
		Username: dto.Login,
		Email:    dto.Login,
	})
	if err != nil {
		s.logger.Error("failed to find user for login",
			zap.String("login", dto.Login),
			zap.Error(err),
		)
		return "", fmt.Errorf("error handling find one user query: %w", err)
	}
	if user == nil {
		s.logger.Warn("login attempt for non-existent user",
			zap.String("login", dto.Login),
		)
		return "", &appErrors.ApplicationError{
			Message: "incorrect login or password",
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password))
	if err != nil {
		s.logger.Warn("incorrect password attempt",
			zap.String("login", dto.Login),
			zap.String("user_id", user.Id),
		)
		return "", &appErrors.ApplicationError{
			Message: "incorrect login or password",
		}
	}

	token, err := s.tokenGenerator.Generate(user)
	if err != nil {
		s.logger.Error("failed to generate token",
			zap.String("user_id", user.Id),
			zap.Error(err),
		)
		return "", fmt.Errorf("error generating token for user: %w", err)
	}

	s.logger.Info("user logged in",
		zap.String("user_id", user.Id),
		zap.String("username", user.Username),
	)

	return token, nil
}
