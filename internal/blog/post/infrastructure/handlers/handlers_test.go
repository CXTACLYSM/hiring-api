package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/mocks"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/services"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"
	entities2 "github.com/CXTACLYSM/hiring-api/pkg/shared/domain"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/middlewares"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/validation"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

const (
	UserId = "550e8400-e29b-41d4-a716-446655440000"
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

func TestFindOnePostHandler_ServeHTTP_Success(t *testing.T) {
	f := setupHTTPTest(t)

	post := entities.NewPost("post test", "post test content", UserId)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+post.Id, nil)
	req.SetPathValue("id", post.Id)
	req.Header.Set("Content-Type", "application/json")

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", post.Id)

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
	ctx = context.WithValue(ctx, middlewares.UserKey, &entities2.User{
		Id:       UserId,
		Username: "test username",
		Email:    "test email",
	})
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	f.findOnePost.EXPECT().
		Handle(findOne.Query{
			Id:     post.Id,
			UserId: UserId,
		}).
		Return(post, nil)

	f.findOnePostHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
