package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/commands/createOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/mocks"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/queries/findList"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/services"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/responses"
	pkgEntities "github.com/CXTACLYSM/hiring-api/pkg/shared/domain"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/middlewares"
	pkgResponses "github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/responses"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/validation"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type httpFixtures struct {
	ctrl *gomock.Controller

	findOnePost   *mocks.MockFindOnePostHandler
	findListPost  *mocks.MockFindListPostHandler
	createOnePost *mocks.MockCreateOnePostHandler
	updateOnePost *mocks.MockUpdateOnePostHandler
	deleteOnePost *mocks.MockDeleteOnePostHandler

	findOnePostHandler   *FindOnePostHandler
	findListPostHandler  *FindListPostHandler
	createOnePostHandler *CreateOnePostHandler
	updateOnePostHandler *UpdateOnePostHandler
	deleteOnePostHandler *DeleteOnePostHandler
}

func setupHTTPTest(t *testing.T) *httpFixtures {
	ctrl := gomock.NewController(t)

	findOnePost := mocks.NewMockFindOnePostHandler(ctrl)
	findListPost := mocks.NewMockFindListPostHandler(ctrl)
	createOnePost := mocks.NewMockCreateOnePostHandler(ctrl)
	updateOnePost := mocks.NewMockUpdateOnePostHandler(ctrl)
	deleteOnePost := mocks.NewMockDeleteOnePostHandler(ctrl)

	v := validation.NewValidator(validator.New())
	logger := zap.NewNop()
	postService := services.NewPostService(
		v,
		findOnePost,
		findListPost,
		createOnePost,
		updateOnePost,
		deleteOnePost,
		logger,
	)

	return &httpFixtures{
		ctrl: ctrl,

		findOnePost:   findOnePost,
		findListPost:  findListPost,
		createOnePost: createOnePost,
		updateOnePost: updateOnePost,
		deleteOnePost: deleteOnePost,

		findOnePostHandler:   NewFindOnePostHandler(postService),
		findListPostHandler:  NewFindListPostHandler(postService),
		createOnePostHandler: NewCreateOnePostHandler(postService),
		updateOnePostHandler: NewUpdateOnePostHandler(postService),
		deleteOnePostHandler: NewDeleteOnePostHandler(postService),
	}
}

func testUser() *pkgEntities.User {
	return &pkgEntities.User{
		Id:       "550e8400-e29b-41d4-a716-446655440000",
		Username: "test username",
		Email:    "test email",
	}
}

func jsonBody(t *testing.T, v any) *bytes.Reader {
	t.Helper()
	data, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewReader(data)
}

func TestFindOnePostHandler_ServeHTTP_Success(t *testing.T) {
	f := setupHTTPTest(t)

	user := testUser()
	post := entities.NewPost("post test", "post test content", user.Id)
	expectedResponse := responses.NewFindOneResponse(post)

	f.findOnePost.EXPECT().
		Handle(findOne.Query{
			Id:     post.Id,
			UserId: user.Id,
		}).
		Return(post, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+post.Id, nil)
	req.Header.Set("Content-Type", "application/json")

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", post.Id)

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
	ctx = context.WithValue(ctx, middlewares.UserKey, user)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	f.findOnePostHandler.ServeHTTP(rec, req)

	var response *responses.FindOneResponse
	json.NewDecoder(rec.Body).Decode(&response)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedResponse, response)
}

func TestFindOnePostHandler_ServeHTTP_UserNotFound(t *testing.T) {
	f := setupHTTPTest(t)

	user := testUser()
	post := entities.NewPost("test name", "test content", user.Id)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+post.Id, nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	f.findOnePostHandler.ServeHTTP(rec, req)

	var response pkgResponses.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&response)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "unauthorized", response.Message)
}

func TestFindOnePostHandler_ServeHTTP_ServiceError(t *testing.T) {
	f := setupHTTPTest(t)

	user := testUser()
	post := entities.NewPost("test name", "test content", user.Id)

	f.findOnePost.EXPECT().
		Handle(findOne.Query{
			UserId: user.Id,
			Id:     post.Id,
		}).
		Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+post.Id, nil)
	req.Header.Set("Content-Type", "application/json")

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", post.Id)

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
	ctx = context.WithValue(ctx, middlewares.UserKey, user)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	f.findOnePostHandler.ServeHTTP(rec, req)

	var response pkgResponses.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&response)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "internal server error", response.Message)
}

func TestFindListPostHandler_ServeHTTP_Success(t *testing.T) {
	f := setupHTTPTest(t)

	user := testUser()
	post := entities.NewPost("test_name", "test_content", user.Id)
	post2 := entities.NewPost("test_name 2", "test_content 2", user.Id)
	filter := map[string]string{
		"name":    "test_name",
		"content": "test_content",
	}
	expectedResponse := []*responses.FindOneResponse{
		responses.NewFindOneResponse(post),
		responses.NewFindOneResponse(post2),
	}

	f.findListPost.EXPECT().
		Handle(findList.Query{
			UserId:  user.Id,
			Name:    filter["name"],
			Content: filter["content"],
		}).
		Return([]*entities.Post{post, post2}, nil)

	path := fmt.Sprintf("/api/v1/posts?name=%s&content=%s", filter["name"], filter["content"])
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), middlewares.UserKey, user)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	f.findListPostHandler.ServeHTTP(rec, req)

	var response []*responses.FindOneResponse
	json.NewDecoder(rec.Body).Decode(&response)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedResponse, response)
}

func TestFindListPostHandler_ServeHTTP_UserNotFound(t *testing.T) {
	f := setupHTTPTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)
	rec := httptest.NewRecorder()

	f.findListPostHandler.ServeHTTP(rec, req)

	var response pkgResponses.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&response)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "unauthorized", response.Message)
}

func TestFindListPostHandler_ServeHTTP_ServiceError(t *testing.T) {
	f := setupHTTPTest(t)

	user := testUser()
	filter := map[string]string{
		"name":    "test_name",
		"content": "test_content",
	}
	f.findListPost.EXPECT().
		Handle(findList.Query{
			UserId:  user.Id,
			Name:    filter["name"],
			Content: filter["content"],
		}).
		Return(nil, assert.AnError)

	path := fmt.Sprintf("/api/v1/posts?name=%s&content=%s", filter["name"], filter["content"])
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), middlewares.UserKey, user)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	f.findListPostHandler.ServeHTTP(rec, req)

	var response pkgResponses.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&response)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "internal server error", response.Message)
}

func TestCreateOnePostHandler_ServeHTTP_Success(t *testing.T) {
	f := setupHTTPTest(t)

	data := map[string]string{
		"name":    "test name",
		"content": "test content",
	}
	user := testUser()
	post := entities.NewPost(data["name"], data["content"], user.Id)
	body := jsonBody(t, data)
	expectedResponse := responses.CreateOneResponse{
		Id:      post.Id,
		Name:    post.Name,
		Content: post.Content,
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", body)
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), middlewares.UserKey, user)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	f.createOnePost.EXPECT().
		Handle(createOne.Command{
			UserId:  user.Id,
			Name:    data["name"],
			Content: data["content"],
		}).
		Return(post, nil)

	f.createOnePostHandler.ServeHTTP(rec, req)

	var response responses.CreateOneResponse
	json.NewDecoder(rec.Body).Decode(&response)
	require.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, expectedResponse, response)
}

func TestCreateOnePostHandler_ServeHTTP_EmptyBody(t *testing.T) {
	f := setupHTTPTest(t)

	user := testUser()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), middlewares.UserKey, user)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	f.createOnePostHandler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var errorResponse pkgResponses.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errorResponse)
	assert.Equal(t, "request body is required", errorResponse.Message)
}

func TestCreateOnePostHandler_ServeHTTP_InvalidJSON(t *testing.T) {
	f := setupHTTPTest(t)

	user := testUser()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader([]byte("{broken")))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), middlewares.UserKey, user)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	f.createOnePostHandler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var errorResponse pkgResponses.ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errorResponse)
	assert.Equal(t, "invalid json syntax", errorResponse.Message)
}

func TestCreateOnePostHandler_ServeHTTP_ValidationFails(t *testing.T) {
	f := setupHTTPTest(t)

	user := testUser()
	body := jsonBody(t, map[string]string{
		"name":    " ",
		"content": " ",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", body)
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), middlewares.UserKey, user)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	f.createOnePostHandler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateOnePostHandler_ServeHTTP_ServiceFail(t *testing.T) {
	// TODO
}

func TestUpdateOnePostHandler_ServeHTTP_Success(t *testing.T) {
	// TODO
}

func TestUpdateOnePostHandler_ServeHTTP_EmptyBody(t *testing.T) {
	// TODO
}

func TestUpdateOnePostHandler_ServeHTTP_InvalidJSON(t *testing.T) {
	// TODO
}

func TestUpdateOnePostHandler_ServeHTTP_ValidationFails(t *testing.T) {
	// TODO
}

func TestUpdateOnePostHandler_ServeHTTP_ServiceFail(t *testing.T) {
	// TODO
}

func TestDeleteOnePostHandler_ServeHTTP_Success(t *testing.T) {
	// TODO
}

func TestDeleteOnePostHandler_ServeHTTP_ServiceFail(t *testing.T) {
	// TODO
}
