package handlers

import (
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/dto"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/services"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/responses"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/httputils"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/middlewares"
	"github.com/go-chi/chi/v5"
)

type UpdateOnePostHandler struct {
	postService *services.PostService
}

func NewUpdateOnePostHandler(postService *services.PostService) *UpdateOnePostHandler {
	return &UpdateOnePostHandler{
		postService: postService,
	}
}

// ServeHTTP godoc
// @Summary     Update post
// @Description Updates an existing post by ID
// @Tags        posts
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id      path     string true "Post ID (UUID)"
// @Param       request body     dto.UpdateOneDTO true "Updated post data"
// @Success     200 {object}     responses.UpdateOneResponse
// @Failure     401 {object}     responses.ErrorResponse "Unauthorized"
// @Failure     422 {object}     responses.ValidationErrorsResponse "Validation errors"
// @Failure     500 {object}     responses.ErrorResponse "Internal server error"
// @Router      /posts/{id} [put]
func (h *UpdateOnePostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserFromRequest(r)
	if user == nil {
		httputils.ResponseError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var updateOneDTO dto.UpdateOneDTO
	if err := httputils.DecodeJSON(r, &updateOneDTO); err != nil {
		httputils.WriteError(w, err)
		return
	}

	updateOneDTO.UserId = user.Id
	updateOneDTO.Id = chi.URLParam(r, "id")

	post, err := h.postService.UpdateOne(&updateOneDTO)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}
	if post == nil {
		httputils.ResponseError(w, http.StatusNotFound, "post not found")
		return
	}

	httputils.ResponseOk(w, http.StatusOK, responses.NewUpdateOneResponse(post))
}
