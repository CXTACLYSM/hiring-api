package services

import (
	"testing"
	"time"

	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/commands/createOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/dto"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/mocks"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"
	appErrors "github.com/CXTACLYSM/hiring-api/pkg/shared/application/errors"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/validation"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type testFixtures struct {
	ctrl           *gomock.Controller
	findOneUser    *mocks.MockFindOneUserHandler
	createOneUser  *mocks.MockCreateOneUserHandler
	tokenGenerator *mocks.MockTokenGenerator
	publisher      *mocks.MockEventPublisher
	authService    *AuthService
}

func setupTest(t *testing.T) *testFixtures {
	ctrl := gomock.NewController(t)

	findOneUser := mocks.NewMockFindOneUserHandler(ctrl)
	createOneUser := mocks.NewMockCreateOneUserHandler(ctrl)
	tokenGenerator := mocks.NewMockTokenGenerator(ctrl)
	publisher := mocks.NewMockEventPublisher(ctrl)

	v := validation.NewValidator(validator.New())

	authService := NewAuthService(
		v,
		findOneUser,
		createOneUser,
		tokenGenerator,
		publisher,
		zap.NewNop(),
	)

	return &testFixtures{
		ctrl:           ctrl,
		findOneUser:    findOneUser,
		createOneUser:  createOneUser,
		tokenGenerator: tokenGenerator,
		publisher:      publisher,
		authService:    authService,
	}
}

func validRegisterDTO() *dto.RegisterDTO {
	return &dto.RegisterDTO{
		Username:             "testuser",
		Email:                "test@example.com",
		Password:             "password123",
		PasswordConfirmation: "password123",
	}
}

func testUser() *entities.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	return &entities.User{
		Id:           "550e8400-e29b-41d4-a716-446655440000",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func TestRegister_Success(t *testing.T) {
	f := setupTest(t)

	user := testUser()
	registerDTO := validRegisterDTO()

	f.findOneUser.EXPECT().
		Handle(findOne.Query{
			Username: registerDTO.Username,
			Email:    registerDTO.Email,
		}).
		Return(nil, nil)

	f.createOneUser.EXPECT().
		Handle(gomock.Any()).
		DoAndReturn(func(cmd createOne.Command) (*entities.User, error) {
			assert.Equal(t, registerDTO.Username, cmd.Username)
			assert.Equal(t, registerDTO.Email, cmd.Email)
			assert.NotEmpty(t, cmd.PasswordHash)
			return user, nil
		})

	f.tokenGenerator.EXPECT().
		Generate(user).
		Return("jwt-token-123", nil)

	f.publisher.EXPECT().
		Push(gomock.Any()).
		Return(nil)

	token, err := f.authService.Register(registerDTO)

	require.NoError(t, err)
	assert.Equal(t, "jwt-token-123", token)
}

func TestRegister_ValidationFails(t *testing.T) {
	f := setupTest(t)

	registerDTO := &dto.RegisterDTO{
		Username:             "",
		Email:                "test@example.com",
		Password:             "password123",
		PasswordConfirmation: "password123",
	}

	token, err := f.authService.Register(registerDTO)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "validating")
}

func TestRegister_PasswordMismatch(t *testing.T) {
	f := setupTest(t)

	registerDTO := &dto.RegisterDTO{
		Username:             "testuser",
		Email:                "test@example.com",
		Password:             "password123",
		PasswordConfirmation: "different456",
	}

	token, err := f.authService.Register(registerDTO)

	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	f := setupTest(t)

	existingUser := testUser()
	registerDTO := validRegisterDTO()

	f.findOneUser.EXPECT().
		Handle(findOne.Query{
			Username: registerDTO.Username,
			Email:    registerDTO.Email,
		}).
		Return(existingUser, nil)

	token, err := f.authService.Register(registerDTO)

	assert.Error(t, err)
	assert.Empty(t, token)

	var appErr *appErrors.ApplicationError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, "user already exists", appErr.Message)
}

func TestRegister_FindUserQueryFails(t *testing.T) {
	f := setupTest(t)

	registerDTO := validRegisterDTO()

	f.findOneUser.EXPECT().
		Handle(gomock.Any()).
		Return(nil, assert.AnError)

	token, err := f.authService.Register(registerDTO)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "find one user")
}

func TestRegister_CreateUserFails(t *testing.T) {
	f := setupTest(t)

	registerDTO := validRegisterDTO()

	f.findOneUser.EXPECT().
		Handle(gomock.Any()).
		Return(nil, nil)

	f.createOneUser.EXPECT().
		Handle(gomock.Any()).
		Return(nil, assert.AnError)

	token, err := f.authService.Register(registerDTO)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "create one user")
}

func TestRegister_TokenGenerationFails(t *testing.T) {
	f := setupTest(t)

	user := testUser()
	registerDTO := validRegisterDTO()

	f.findOneUser.EXPECT().Handle(gomock.Any()).Return(nil, nil)
	f.createOneUser.EXPECT().Handle(gomock.Any()).Return(user, nil)
	f.tokenGenerator.EXPECT().Generate(user).Return("", assert.AnError)

	token, err := f.authService.Register(registerDTO)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "generating token")
}

func TestRegister_KafkaFailsButRegistrationSucceeds(t *testing.T) {
	f := setupTest(t)

	user := testUser()
	registerDTO := validRegisterDTO()

	f.findOneUser.EXPECT().Handle(gomock.Any()).Return(nil, nil)
	f.createOneUser.EXPECT().Handle(gomock.Any()).Return(user, nil)
	f.tokenGenerator.EXPECT().Generate(user).Return("jwt-token-123", nil)

	f.publisher.EXPECT().
		Push(gomock.Any()).
		Return(assert.AnError)

	token, err := f.authService.Register(registerDTO)

	require.NoError(t, err)
	assert.Equal(t, "jwt-token-123", token)
}

func TestLogin_Success(t *testing.T) {
	f := setupTest(t)

	user := testUser()

	loginDTO := &dto.LoginDTO{
		Login:    "testuser",
		Password: "password123",
	}

	f.findOneUser.EXPECT().
		Handle(findOne.Query{
			Username: loginDTO.Login,
			Email:    loginDTO.Login,
		}).
		Return(user, nil)

	f.tokenGenerator.EXPECT().
		Generate(user).
		Return("jwt-token-456", nil)

	token, err := f.authService.Login(loginDTO)

	require.NoError(t, err)
	assert.Equal(t, "jwt-token-456", token)
}

func TestLogin_ValidationFails(t *testing.T) {
	f := setupTest(t)

	loginDTO := &dto.LoginDTO{
		Login:    "",
		Password: "password123",
	}

	token, err := f.authService.Login(loginDTO)

	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestLogin_UserNotFound(t *testing.T) {
	f := setupTest(t)

	loginDTO := &dto.LoginDTO{
		Login:    "nonexistent",
		Password: "password123",
	}

	f.findOneUser.EXPECT().
		Handle(gomock.Any()).
		Return(nil, nil)

	token, err := f.authService.Login(loginDTO)

	assert.Error(t, err)
	assert.Empty(t, token)

	var appErr *appErrors.ApplicationError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, "incorrect login or password", appErr.Message)
}

func TestLogin_WrongPassword(t *testing.T) {
	f := setupTest(t)

	user := testUser()

	loginDTO := &dto.LoginDTO{
		Login:    "testuser",
		Password: "wrongpassword",
	}

	f.findOneUser.EXPECT().
		Handle(gomock.Any()).
		Return(user, nil)

	token, err := f.authService.Login(loginDTO)

	assert.Error(t, err)
	assert.Empty(t, token)

	var appErr *appErrors.ApplicationError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, "incorrect login or password", appErr.Message)
}

func TestLogin_FindUserQueryFails(t *testing.T) {
	f := setupTest(t)

	loginDTO := &dto.LoginDTO{
		Login:    "testuser",
		Password: "password123",
	}

	f.findOneUser.EXPECT().
		Handle(gomock.Any()).
		Return(nil, assert.AnError)

	token, err := f.authService.Login(loginDTO)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "find one user")
}

func TestLogin_TokenGenerationFails(t *testing.T) {
	f := setupTest(t)

	user := testUser()

	loginDTO := &dto.LoginDTO{
		Login:    "testuser",
		Password: "password123",
	}

	f.findOneUser.EXPECT().Handle(gomock.Any()).Return(user, nil)
	f.tokenGenerator.EXPECT().Generate(user).Return("", assert.AnError)

	token, err := f.authService.Login(loginDTO)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "generating token")
}
