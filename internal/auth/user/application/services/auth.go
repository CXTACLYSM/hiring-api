package services

import (
	"fmt"

	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/commands/createOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/dto"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/queries/findOne"
	appErrors "github.com/CXTACLYSM/hiring-api/pkg/shared/application/errors"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/validation"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	validator      *validation.Validator
	findOneUser    findOne.Handler
	createOneUser  createOne.Handler
	tokenGenerator TokenGenerator
}

func NewAuthService(validator *validation.Validator, findOneUser findOne.Handler, createOneUser createOne.Handler, tokenGenerator TokenGenerator) *AuthService {
	return &AuthService{
		validator:      validator,
		findOneUser:    findOneUser,
		createOneUser:  createOneUser,
		tokenGenerator: tokenGenerator,
	}
}

func (s *AuthService) Register(dto *dto.RegisterDTO) (string, error) {
	if err := s.validator.Struct(dto); err != nil {
		return "", fmt.Errorf("error validating register dto: %w", err)
	}

	user, err := s.findOneUser.Handle(
		findOne.Query{
			Username: dto.Username,
			Email:    dto.Email,
		},
	)
	if err != nil {
		return "", fmt.Errorf("error handling find one user query: %w", err)
	}
	if user != nil {
		return "", &appErrors.ApplicationError{
			Message: "user already exists",
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error generating password hash: %w", err)
	}
	user, err = s.createOneUser.Handle(createOne.Command{
		Username:     dto.Username,
		Email:        dto.Email,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		return "", fmt.Errorf("error handling create one user query: %w", err)
	}

	token, err := s.tokenGenerator.Generate(user)
	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}

	return token, nil
}

func (s *AuthService) Login(dto *dto.LoginDTO) (string, error) {
	if err := s.validator.Struct(dto); err != nil {
		return "", fmt.Errorf("error validation login dto: %w", err)
	}

	user, err := s.findOneUser.Handle(
		findOne.Query{
			Username: dto.Login,
			Email:    dto.Login,
		},
	)
	if err != nil {
		return "", fmt.Errorf("error handling find one user query: %w", err)
	}
	if user == nil {
		return "", &appErrors.ApplicationError{
			Message: "incorrect login or password",
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password))
	if err != nil {
		return "", &appErrors.ApplicationError{
			Message: "incorrect login or password",
		}
	}

	token, err := s.tokenGenerator.Generate(user)
	if err != nil {
		return "", fmt.Errorf("error generating token for user: %w", err)
	}

	return token, nil
}
