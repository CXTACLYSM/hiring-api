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

type FindOnePostHandler struct {
	postService *services.PostService
}

func NewFindOnePostHandler(postService *services.PostService) *FindOnePostHandler {
	return &FindOnePostHandler{
		postService: postService,
	}
}

// ServeHTTP godoc
// @Summary     Get post
// @Description Returns a single post by ID
// @Tags        posts
// @Produce     json
// @Security    BearerAuth
// @Param       id   path       string true "Post ID (UUID)"
// @Success     200 {object}    responses.FindOneResponse
// @Failure     401 {object}    responses.ErrorResponse "Unauthorized"
// @Failure     404 {object}    responses.ErrorResponse "Post not found"
// @Failure     500 {object}    responses.ErrorResponse "Internal server error"
// @Router      /posts/{id} [get]
func (h *FindOnePostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserFromRequest(r)
	if user == nil {
		httputils.ResponseError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	findOneDTO := &dto.FindOneDTO{
		Id:     chi.URLParam(r, "id"),
		UserId: user.Id,
	}

	post, err := h.postService.FindOne(findOneDTO)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}
	if post == nil {
		httputils.ResponseError(w, http.StatusNotFound, "post not found")
		return
	}

	httputils.ResponseOk(w, http.StatusOK, responses.NewFindOneResponse(post))
}
