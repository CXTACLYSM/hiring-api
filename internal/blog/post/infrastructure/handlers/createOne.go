package handlers

import (
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/dto"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/services"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/responses"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/httputils"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/middlewares"
)

type CreateOnePostHandler struct {
	postService *services.PostService
}

func NewCreateOnePostHandler(postService *services.PostService) *CreateOnePostHandler {
	return &CreateOnePostHandler{
		postService: postService,
	}
}

// ServeHTTP godoc
// @Summary     Create post
// @Description Creates a new blog post for the authenticated user
// @Tags        posts
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body     dto.CreateOneDTO true "Post data"
// @Success     201 {object}     responses.CreateOneResponse
// @Failure     401 {object}     responses.ErrorResponse "Unauthorized"
// @Failure     422 {object}     responses.ValidationErrorsResponse "Validation errors"
// @Failure     500 {object}     responses.ErrorResponse "Internal server error"
// @Router      /posts [post]
func (h *CreateOnePostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserFromRequest(r)
	if user == nil {
		httputils.ResponseError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var createOneDTO dto.CreateOneDTO
	if err := httputils.DecodeJSON(r, &createOneDTO); err != nil {
		httputils.WriteError(w, err)
		return
	}
	createOneDTO.UserId = user.Id

	post, err := h.postService.CreateOne(&createOneDTO)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.ResponseOk(w, http.StatusCreated, responses.NewCreateOneResponse(post))
}
