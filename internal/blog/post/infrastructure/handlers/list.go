package handlers

import (
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/dto"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/application/services"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/responses"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/httputils"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/middlewares"
)

type FindListPostHandler struct {
	postService *services.PostService
}

func NewFindListPostHandler(postService *services.PostService) *FindListPostHandler {
	return &FindListPostHandler{
		postService: postService,
	}
}

func (h *FindListPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserFromRequest(r)
	if user == nil {
		httputils.ResponseError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	listDTO := &dto.ListDTO{
		Name:    r.URL.Query().Get("name"),
		Content: r.URL.Query().Get("content"),
		UserId:  user.Id,
	}

	postList, err := h.postService.FindList(listDTO)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.ResponseOk(w, http.StatusOK, responses.NewFindListResponse(postList))
}
