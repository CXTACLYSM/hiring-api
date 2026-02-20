package handlers

import (
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/dto"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/services"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/httputils"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/middlewares"
	"github.com/go-chi/chi/v5"
)

type DeleteOnePostHandler struct {
	postService *services.PostService
}

func NewDeleteOnePostHandler(postService *services.PostService) *DeleteOnePostHandler {
	return &DeleteOnePostHandler{
		postService: postService,
	}
}

func (h *DeleteOnePostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserFromRequest(r)
	if user == nil {
		httputils.ResponseError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	deleteOneDTO := &dto.DeleteOneDTO{
		Id:     chi.URLParam(r, "id"),
		UserId: user.Id,
	}

	if err := h.postService.DeleteOne(deleteOneDTO); err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.ResponseOkNoData(w, http.StatusAccepted)
}
