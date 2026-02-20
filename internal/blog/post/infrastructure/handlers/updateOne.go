package handlers

import (
	"errors"
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

func (h *UpdateOnePostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserFromRequest(r)
	if user == nil {
		httputils.ResponseError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var updateOneDTO *dto.UpdateOneDTO
	if err := httputils.DecodeJSON(r, &updateOneDTO); err != nil {
		httputils.WriteError(w, err)
		return
	}

	updateOneDTO.UserId = user.Id
	updateOneDTO.Id = chi.URLParam(r, "id")

	user = middlewares.UserFromRequest(r)
	if user == nil {
		httputils.WriteError(w, errors.New("claims is nil"))
		return
	}
	updateOneDTO.UserId = user.Id

	post, err := h.postService.UpdateOne(updateOneDTO)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.ResponseOk(w, http.StatusOK, responses.NewUpdateOneResponse(post))
}
