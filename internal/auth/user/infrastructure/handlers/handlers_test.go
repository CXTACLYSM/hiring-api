package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/mocks"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/services"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/infrastructure/responses"
	entities2 "github.com/CXTACLYSM/hiring-api/pkg/shared/domain"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/middlewares"
	pkgResponses "github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/responses"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/validation"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type httpFixtures struct {
	ctrl *gomock.Controller

	findOneUser    *mocks.MockFindOneUserHandler
	createOneUser  *mocks.MockCreateOneUserHandler
	tokenGenerator *mocks.MockTokenGenerator
	publisher      *mocks.MockEventPublisher

	registerHandler *RegisterHandler
	loginHandler    *LoginHandler
	meHandler       *MeHandler
}

func setupHTTPTest(t *testing.T) *httpFixtures {
	ctrl := gomock.NewController(t)

	findOneUser := mocks.NewMockFindOneUserHandler(ctrl)
	createOneUser := mocks.NewMockCreateOneUserHandler(ctrl)
	tokenGenerator := mocks.NewMockTokenGenerator(ctrl)
	publisher := mocks.NewMockEventPublisher(ctrl)

	v := validation.NewValidator(validator.New())

	logger := zap.NewNop()
	authService := services.NewAuthService(
		v,
		findOneUser,
		createOneUser,
		tokenGenerator,
		publisher,
		logger,
	)
	userService := services.NewUserService(
		findOneUser,
		createOneUser,
		logger,
	)

	return &httpFixtures{
		ctrl: ctrl,

		findOneUser:    findOneUser,
		createOneUser:  createOneUser,
		tokenGenerator: tokenGenerator,
		publisher:      publisher,

		registerHandler: NewRegisterHandler(authService),
		loginHandler:    NewLoginHandler(authService),
		meHandler:       NewMeHandler(userService),
	}
}

func jsonBody(t *testing.T, v any) *bytes.Reader {
	t.Helper()
	data, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewReader(data)
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

func TestRegisterHandler_Success(t *testing.T) {
	f := setupHTTPTest(t)

	user := testUser()

	f.findOneUser.EXPECT().Handle(gomock.Any()).Return(nil, nil)
	f.createOneUser.EXPECT().Handle(gomock.Any()).Return(user, nil)
	f.tokenGenerator.EXPECT().Generate(user).Return("jwt-token-123", nil)
	f.publisher.EXPECT().Push(gomock.Any()).Return(nil)

	body := jsonBody(t, map[string]string{
		"username":              "testuser",
		"email":                 "test@example.com",
		"password":              "password123",
		"password_confirmation": "password123",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	f.registerHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp responses.RegisterResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "jwt-token-123", resp.Token)
}

func TestRegisterHandler_EmptyBody(t *testing.T) {
	f := setupHTTPTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	f.registerHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp pkgResponses.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "request body is required", resp.Message)
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	f := setupHTTPTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader([]byte("{broken")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	f.registerHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp pkgResponses.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "invalid json syntax", resp.Message)
}

func TestRegisterHandler_ValidationFails(t *testing.T) {
	f := setupHTTPTest(t)

	body := jsonBody(t, map[string]string{
		"username": "",
		"email":    "test@example.com",
		"password": "password123",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	f.registerHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestRegisterHandler_UserAlreadyExists(t *testing.T) {
	f := setupHTTPTest(t)

	existingUser := testUser()

	f.findOneUser.EXPECT().Handle(gomock.Any()).Return(existingUser, nil)

	body := jsonBody(t, map[string]string{
		"username":              "testuser",
		"email":                 "test@example.com",
		"password":              "password123",
		"password_confirmation": "password123",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	f.registerHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp pkgResponses.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "user already exists", resp.Message)
}

func TestRegisterHandler_InternalError(t *testing.T) {
	f := setupHTTPTest(t)

	f.findOneUser.EXPECT().Handle(gomock.Any()).Return(nil, assert.AnError)

	body := jsonBody(t, map[string]string{
		"username":              "testuser",
		"email":                 "test@example.com",
		"password":              "password123",
		"password_confirmation": "password123",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	f.registerHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp pkgResponses.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "internal server error", resp.Message)
}

func TestLoginHandler_Success(t *testing.T) {
	f := setupHTTPTest(t)

	user := testUser()

	f.findOneUser.EXPECT().Handle(gomock.Any()).Return(user, nil)
	f.tokenGenerator.EXPECT().Generate(user).Return("jwt-token-456", nil)

	body := jsonBody(t, map[string]string{
		"login":    "testuser",
		"password": "password123",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	f.loginHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp responses.LoginResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "jwt-token-456", resp.Token)
}

func TestLoginHandler_EmptyBody(t *testing.T) {
	f := setupHTTPTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	f.loginHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp pkgResponses.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "request body is required", resp.Message)
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	f := setupHTTPTest(t)

	user := testUser()

	f.findOneUser.EXPECT().Handle(gomock.Any()).Return(user, nil)

	body := jsonBody(t, map[string]string{
		"login":    "testuser",
		"password": "wrongpassword",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	f.loginHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp pkgResponses.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "incorrect login or password", resp.Message)
}

func TestLoginHandler_UserNotFound(t *testing.T) {
	f := setupHTTPTest(t)

	f.findOneUser.EXPECT().Handle(gomock.Any()).Return(nil, nil)

	body := jsonBody(t, map[string]string{
		"login":    "nonexistent",
		"password": "password123",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	f.loginHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp pkgResponses.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "incorrect login or password", resp.Message)
}

func TestMeHandler_Success(t *testing.T) {
	f := setupHTTPTest(t)

	user := testUser()

	f.findOneUser.EXPECT().
		Handle(findOne.Query{
			Id: user.Id,
		}).
		Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), middlewares.UserKey, &entities2.User{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	})
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	f.meHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp responses.MeResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, user.Id, resp.Id)
}

func TestMeHandler_UserNotFound(t *testing.T) {
	f := setupHTTPTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	f.meHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var resp pkgResponses.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "unauthorized", resp.Message)
}
