package services

import (
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/auth/dto"
	appErrors "github.com/CXTACLYSM/hiring-api/internal/auth/errors"
	"github.com/CXTACLYSM/hiring-api/internal/auth/queries/user/createOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/queries/user/findOne"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	findOneUser    *findOne.Handler
	createOneUser  *createOne.Handler
	tokenGenerator TokenGenerator
}

func NewAuthService(findOne *findOne.Handler, createOne *createOne.Handler, tokenGenerator TokenGenerator) *AuthService {
	return &AuthService{
		findOneUser:    findOne,
		createOneUser:  createOne,
		tokenGenerator: tokenGenerator,
	}
}

func (s *AuthService) Register(dto *dto.RegisterDTO) (string, error) {
	if dto.Username == "" {
		return "", &appErrors.AppError{
			Message:    "fill username",
			StatusCode: http.StatusUnprocessableEntity,
			Err:        nil,
		}
	}
	if dto.Email == "" {
		return "", &appErrors.AppError{
			Message:    "fill email",
			StatusCode: http.StatusUnprocessableEntity,
			Err:        nil,
		}
	}
	if dto.Password == "" {
		return "", &appErrors.AppError{
			Message:    "fill password",
			StatusCode: http.StatusUnprocessableEntity,
			Err:        nil,
		}
	}
	if dto.Password != dto.PasswordConfirmation {
		return "", &appErrors.AppError{
			Message:    "passwords must be same",
			StatusCode: http.StatusUnprocessableEntity,
			Err:        nil,
		}
	}

	user, err := s.findOneUser.Handle(findOne.Query{
		Username: dto.Username,
		Email:    dto.Email,
	})
	if err != nil {
		return "", err
	}
	if user != nil {
		return "", &appErrors.AppError{
			Message:    "user already exists",
			StatusCode: http.StatusUnprocessableEntity,
			Err:        nil,
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user, err = s.createOneUser.Handle(createOne.Query{
		Username:     dto.Username,
		Email:        dto.Email,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		return "", err
	}

	token, err := s.tokenGenerator.Generate(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Login(dto *dto.LoginDTO) (string, error) {
	if len(dto.Login) == 0 {
		return "", &appErrors.AppError{
			Message:    "login required (email of username)",
			StatusCode: http.StatusUnprocessableEntity,
			Err:        nil,
		}
	}
	if len(dto.Password) == 0 {
		return "", &appErrors.AppError{
			Message:    "password required",
			StatusCode: http.StatusUnprocessableEntity,
			Err:        nil,
		}
	}

	user, err := s.findOneUser.Handle(findOne.Query{
		Username: dto.Login,
		Email:    dto.Login,
	})
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", &appErrors.AppError{
			Message:    "incorrect login or password",
			StatusCode: http.StatusUnprocessableEntity,
			Err:        nil,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password))
	if err != nil {
		return "", &appErrors.AppError{
			Message:    "incorrect login or password",
			StatusCode: http.StatusUnprocessableEntity,
			Err:        nil,
		}
	}

	token, err := s.tokenGenerator.Generate(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
