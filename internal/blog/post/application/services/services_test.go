package services

import (
	"testing"

	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/commands/createOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/commands/deleteOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/commands/updateOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/dto"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/mocks"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/queries/findList"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/validation"
	"github.com/CXTACLYSM/hiring-api/pkg/tests"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

const (
	UserId = "550e8400-e29b-41d4-a716-446655440000"
)

type testFixtures struct {
	ctrl *gomock.Controller

	findOnePost   *mocks.MockFindOnePostHandler
	findListPost  *mocks.MockFindListPostHandler
	createOnePost *mocks.MockCreateOnePostHandler
	updateOnePost *mocks.MockUpdateOnePostHandler
	deleteOnePost *mocks.MockDeleteOnePostHandler

	postService *PostService
}

func setupTest(t *testing.T) *testFixtures {
	ctrl := gomock.NewController(t)

	findOnePost := mocks.NewMockFindOnePostHandler(ctrl)
	findListPost := mocks.NewMockFindListPostHandler(ctrl)
	createOnePost := mocks.NewMockCreateOnePostHandler(ctrl)
	updateOnePost := mocks.NewMockUpdateOnePostHandler(ctrl)
	deleteOnePost := mocks.NewMockDeleteOnePostHandler(ctrl)

	v := validation.NewValidator(validator.New())
	logger := zap.NewNop()

	postService := NewPostService(
		v,
		findOnePost,
		findListPost,
		createOnePost,
		updateOnePost,
		deleteOnePost,
		logger,
	)

	return &testFixtures{
		ctrl:          ctrl,
		findOnePost:   findOnePost,
		findListPost:  findListPost,
		createOnePost: createOnePost,
		updateOnePost: updateOnePost,
		deleteOnePost: deleteOnePost,
		postService:   postService,
	}
}

func validCreateOnePostDTO() *dto.CreateOneDTO {
	return &dto.CreateOneDTO{
		Name:    "test post",
		Content: "test post content",
		UserId:  UserId,
	}
}

func validUpdateOnePostDTO(id string) *dto.UpdateOneDTO {
	return &dto.UpdateOneDTO{
		Id:      id,
		Name:    "test post",
		Content: "test post content",
		UserId:  UserId,
	}
}

func validDeleteOnePostDTO(id string) *dto.DeleteOneDTO {
	return &dto.DeleteOneDTO{
		Id:     id,
		UserId: UserId,
	}
}

func validFindOnePostDTO(id string) *dto.FindOneDTO {
	return &dto.FindOneDTO{
		Id:     id,
		UserId: UserId,
	}
}

func validFindListPostDTO(name, content string) *dto.FindListDTO {
	return &dto.FindListDTO{
		Name:    name,
		Content: content,
		UserId:  UserId,
	}
}

func TestPostService_CreateOne_Success(t *testing.T) {
	f := setupTest(t)

	post := entities.NewPost("test post", "test post content", UserId)

	createOnePostDTO := validCreateOnePostDTO()

	f.createOnePost.EXPECT().
		Handle(createOne.Command{
			UserId:  createOnePostDTO.UserId,
			Name:    createOnePostDTO.Name,
			Content: createOnePostDTO.Content,
		}).
		Return(post, nil)

	result, err := f.postService.CreateOne(createOnePostDTO)

	require.NoError(t, err)
	assert.Equal(t, post, result)
}

func TestPostService_CreateOne_ValidationErrorsRequired(t *testing.T) {
	f := setupTest(t)

	createOnePostDTO := &dto.CreateOneDTO{
		Name:    "",
		Content: "",
		UserId:  "",
	}
	expectedErrors := map[string]string{
		"Name":    "required",
		"Content": "required",
		"UserId":  "required",
	}

	result, err := f.postService.CreateOne(createOnePostDTO)
	require.Error(t, err)
	assert.Nil(t, result)
	tests.RequireValidationErrors(t, expectedErrors, err)
}

func TestPostService_CreateOne_ValidationErrors(t *testing.T) {
	f := setupTest(t)

	createOnePostDTO := &dto.CreateOneDTO{
		Name:    " ",
		Content: " ",
		UserId:  " ",
	}
	expectedErrors := map[string]string{
		"Name":    "min",
		"Content": "min",
		"UserId":  "uuid",
	}

	result, err := f.postService.CreateOne(createOnePostDTO)
	require.Error(t, err)
	assert.Nil(t, result)
	tests.RequireValidationErrors(t, expectedErrors, err)
}

func TestPostService_CreateOne_CommandFail(t *testing.T) {
	f := setupTest(t)

	createOnePostDTO := validCreateOnePostDTO()

	f.createOnePost.EXPECT().
		Handle(createOne.Command{
			UserId:  createOnePostDTO.UserId,
			Name:    createOnePostDTO.Name,
			Content: createOnePostDTO.Content,
		}).
		Return(nil, assert.AnError)

	result, err := f.postService.CreateOne(createOnePostDTO)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error creating")
}

func TestPostService_UpdateOne_Success(t *testing.T) {
	f := setupTest(t)

	post := entities.NewPost("test post", "test post content", UserId)
	updateOnePostDTO := validUpdateOnePostDTO(post.Id)

	freshPost := *post
	freshPost.Name = updateOnePostDTO.Name
	freshPost.Content = updateOnePostDTO.Content

	f.updateOnePost.EXPECT().
		Handle(updateOne.Command{
			UserId:  updateOnePostDTO.UserId,
			Id:      updateOnePostDTO.Id,
			Name:    updateOnePostDTO.Name,
			Content: updateOnePostDTO.Content,
		}).
		Return(&freshPost, nil)

	result, err := f.postService.UpdateOne(updateOnePostDTO)

	require.NoError(t, err)
	assert.Equal(t, &freshPost, result)
}

func TestPostService_UpdateOne_ValidationErrorsRequired(t *testing.T) {
	f := setupTest(t)

	updateOnePostDTO := &dto.UpdateOneDTO{
		Id:      "",
		Name:    "",
		Content: "",
		UserId:  "",
	}
	expectedErrors := map[string]string{
		"Id":      "required",
		"Name":    "required",
		"Content": "required",
		"UserId":  "required",
	}

	result, err := f.postService.UpdateOne(updateOnePostDTO)

	tests.RequireValidationErrors(t, expectedErrors, err)
	assert.Nil(t, result)
}

func TestPostService_UpdateOne_ValidationErrors(t *testing.T) {
	f := setupTest(t)

	updateOnePostDTO := &dto.UpdateOneDTO{
		Id:      " ",
		Name:    " ",
		Content: " ",
		UserId:  " ",
	}
	expectedErrors := map[string]string{
		"Id":      "uuid",
		"Name":    "min",
		"Content": "min",
		"UserId":  "uuid",
	}

	result, err := f.postService.UpdateOne(updateOnePostDTO)

	tests.RequireValidationErrors(t, expectedErrors, err)
	assert.Nil(t, result)
}

func TestPostService_UpdateOne_CommandFail(t *testing.T) {
	f := setupTest(t)

	post := entities.NewPost("test post", "test post content", UserId)
	updateOnePostDTO := validUpdateOnePostDTO(post.Id)

	f.updateOnePost.EXPECT().
		Handle(updateOne.Command{
			UserId:  updateOnePostDTO.UserId,
			Id:      updateOnePostDTO.Id,
			Name:    updateOnePostDTO.Name,
			Content: updateOnePostDTO.Content,
		}).
		Return(nil, assert.AnError)

	result, err := f.postService.UpdateOne(updateOnePostDTO)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error updating")
}

func TestPostService_DeleteOne_Success(t *testing.T) {
	f := setupTest(t)

	post := entities.NewPost("test post", "test post content", UserId)
	deleteOnePostDTO := validDeleteOnePostDTO(post.Id)

	f.deleteOnePost.EXPECT().
		Handle(deleteOne.Command{
			UserId: deleteOnePostDTO.UserId,
			Id:     deleteOnePostDTO.Id,
		}).
		Return(nil)

	err := f.postService.DeleteOne(deleteOnePostDTO)

	require.NoError(t, err)
}

func TestPostService_DeleteOne_ValidationErrorsRequired(t *testing.T) {
	f := setupTest(t)

	deleteOnePostDTO := &dto.DeleteOneDTO{
		Id:     "",
		UserId: "",
	}
	expectedErrors := map[string]string{
		"Id":     "required",
		"UserId": "required",
	}

	err := f.postService.DeleteOne(deleteOnePostDTO)

	tests.RequireValidationErrors(t, expectedErrors, err)
}

func TestPostService_DeleteOne_ValidationErrors(t *testing.T) {
	f := setupTest(t)

	deleteOnePostDTO := &dto.DeleteOneDTO{
		Id:     " ",
		UserId: " ",
	}
	expectedErrors := map[string]string{
		"Id":     "uuid",
		"UserId": "uuid",
	}

	err := f.postService.DeleteOne(deleteOnePostDTO)

	tests.RequireValidationErrors(t, expectedErrors, err)
}

func TestPostService_DeleteOne_CommandFail(t *testing.T) {
	f := setupTest(t)

	post := entities.NewPost("test post", "test post content", UserId)
	deleteOnePostDTO := validDeleteOnePostDTO(post.Id)

	f.deleteOnePost.EXPECT().
		Handle(gomock.Any()).
		Return(assert.AnError)

	err := f.postService.DeleteOne(deleteOnePostDTO)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting")
}

func TestPostService_FindOne_Success(t *testing.T) {
	f := setupTest(t)

	post := entities.NewPost("test post", "test post content", UserId)
	findOnePostDTO := validFindOnePostDTO(post.Id)

	f.findOnePost.EXPECT().
		Handle(findOne.Query{
			UserId: findOnePostDTO.UserId,
			Id:     findOnePostDTO.Id,
		}).
		Return(post, nil)

	result, err := f.postService.FindOne(findOnePostDTO)
	require.NoError(t, err)
	assert.Equal(t, post, result)
}

func TestPostService_FindOne_ValidationErrorsRequired(t *testing.T) {
	f := setupTest(t)

	findOnePostDTO := &dto.FindOneDTO{
		Id:     "",
		UserId: "",
	}
	expectedErrors := map[string]string{
		"Id":     "required",
		"UserId": "required",
	}

	result, err := f.postService.FindOne(findOnePostDTO)

	tests.RequireValidationErrors(t, expectedErrors, err)
	assert.Nil(t, result)
}

func TestPostService_FindOne_ValidationErrors(t *testing.T) {
	f := setupTest(t)

	findOnePostDTO := &dto.FindOneDTO{
		Id:     " ",
		UserId: " ",
	}
	expectedErrors := map[string]string{
		"Id":     "uuid",
		"UserId": "uuid",
	}

	result, err := f.postService.FindOne(findOnePostDTO)

	tests.RequireValidationErrors(t, expectedErrors, err)
	assert.Nil(t, result)
}

func TestPostService_FindOne_QueryFail(t *testing.T) {
	f := setupTest(t)

	post := entities.NewPost("test post", "test post content", UserId)
	findOnePostDTO := validFindOnePostDTO(post.Id)

	f.findOnePost.EXPECT().
		Handle(findOne.Query{
			UserId: findOnePostDTO.UserId,
			Id:     findOnePostDTO.Id,
		}).
		Return(nil, assert.AnError)

	result, err := f.postService.FindOne(findOnePostDTO)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error finding one")
	assert.Nil(t, result)
}

func TestPostService_FindList_Success(t *testing.T) {
	f := setupTest(t)

	post := entities.NewPost("test post", "test post content", UserId)
	post2 := entities.NewPost("test post 2", "test post 2 content", UserId)

	findListPostDTO := validFindListPostDTO("", "")

	f.findListPost.EXPECT().
		Handle(findList.Query{
			UserId:  findListPostDTO.UserId,
			Name:    findListPostDTO.Name,
			Content: findListPostDTO.Content,
		}).
		Return([]*entities.Post{post, post2}, nil)

	result, err := f.postService.FindList(findListPostDTO)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, post, result[0])
	assert.Equal(t, post2, result[1])
}

func TestPostService_FindList_ValidationErrorsRequired(t *testing.T) {
	f := setupTest(t)

	findListPostDTO := &dto.FindListDTO{
		Name:    "",
		Content: "",
		UserId:  "",
	}
	expectedErrors := map[string]string{
		"UserId": "required",
	}

	result, err := f.postService.FindList(findListPostDTO)

	tests.RequireValidationErrors(t, expectedErrors, err)
	assert.Nil(t, result)
}

func TestPostService_FindList_ValidationErrors(t *testing.T) {
	f := setupTest(t)

	findListPostDTO := &dto.FindListDTO{
		Name:    " ",
		Content: " ",
		UserId:  " ",
	}
	expectedErrors := map[string]string{
		"Name":    "min",
		"Content": "min",
		"UserId":  "uuid",
	}

	result, err := f.postService.FindList(findListPostDTO)

	tests.RequireValidationErrors(t, expectedErrors, err)
	assert.Nil(t, result)
}

func TestPostService_FindList_QueryFail(t *testing.T) {
	f := setupTest(t)

	findListPostDTO := validFindListPostDTO("", "")

	f.findListPost.EXPECT().
		Handle(findList.Query{
			UserId:  findListPostDTO.UserId,
			Name:    findListPostDTO.Name,
			Content: findListPostDTO.Content,
		}).
		Return(nil, assert.AnError)

	result, err := f.postService.FindList(findListPostDTO)

	require.Error(t, err)
	assert.Nil(t, result)
}
