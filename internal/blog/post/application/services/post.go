package services

import (
	"fmt"

	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/commands/createOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/commands/deleteOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/commands/updateOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/dto"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/queries/findList"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/domain/entities"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/validation"
)

type PostService struct {
	validator     *validation.Validator
	findOnePost   findOne.Handler
	findListPost  findList.Handler
	createOnePost createOne.Handler
	updateOnePost updateOne.Handler
	deleteOnePost deleteOne.Handler
}

func NewPostService(
	validator *validation.Validator,
	findOnePost findOne.Handler,
	findListPost findList.Handler,
	createOnePost createOne.Handler,
	updateOnePost updateOne.Handler,
	deleteOnePost deleteOne.Handler,
) *PostService {
	return &PostService{
		validator:     validator,
		findOnePost:   findOnePost,
		findListPost:  findListPost,
		createOnePost: createOnePost,
		updateOnePost: updateOnePost,
		deleteOnePost: deleteOnePost,
	}
}

func (s *PostService) CreateOne(dto *dto.CreateOneDTO) (*entities.Post, error) {
	if err := s.validator.Struct(dto); err != nil {
		return nil, err
	}

	post, err := s.createOnePost.Handle(
		createOne.Command{
			Name:    dto.Name,
			Content: dto.Content,
			UserId:  dto.UserId,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error creating one post: %w", err)
	}

	return post, nil
}

func (s *PostService) UpdateOne(dto *dto.UpdateOneDTO) (*entities.Post, error) {
	if err := s.validator.Struct(dto); err != nil {
		return nil, err
	}

	post, err := s.updateOnePost.Handle(
		updateOne.Command{
			Id:      dto.Id,
			Name:    dto.Name,
			Content: dto.Content,
			UserId:  dto.UserId,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error updating one post: %w", err)
	}

	return post, nil
}

func (s *PostService) DeleteOne(dto *dto.DeleteOneDTO) error {
	if err := s.validator.Struct(dto); err != nil {
		return err
	}

	err := s.deleteOnePost.Handle(
		deleteOne.Command{
			Id: dto.Id,
		},
	)
	if err != nil {
		return fmt.Errorf("error deleting one post: %w", err)
	}

	return nil
}

func (s *PostService) FindOne(dto *dto.FindOneDTO) (*entities.Post, error) {
	if err := s.validator.Struct(dto); err != nil {
		return nil, err
	}

	post, err := s.findOnePost.Handle(
		findOne.Query{
			Id: dto.Id,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error finding one post: %w", err)
	}

	return post, nil
}

func (s *PostService) FindList(dto *dto.ListDTO) ([]*entities.Post, error) {
	if err := s.validator.Struct(dto); err != nil {
		return nil, err
	}

	postList, err := s.findListPost.Handle(
		findList.Query{
			Name:    dto.Name,
			Content: dto.Content,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error finding post list: %w", err)
	}

	return postList, nil
}
