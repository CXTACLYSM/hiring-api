package handlers

import (
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/dto"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/services"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/responses"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/httputils"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/middlewares"
)

type FindOnePostHandler struct {
	postService *services.PostService
}

func NewFindOnePostHandler(postService *services.PostService) *FindOnePostHandler {
	return &FindOnePostHandler{
		postService: postService,
	}
}

func (h *FindOnePostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserFromRequest(r)
	if user == nil {
		httputils.ResponseError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	findOneDTO := &dto.FindOneDTO{
		Id:     r.PathValue("id"),
		UserId: user.Id,
	}

	post, err := h.postService.FindOne(findOneDTO)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.ResponseOk(w, http.StatusOK, responses.NewFindOneResponse(post))
}
