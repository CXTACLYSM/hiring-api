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

// ServeHTTP godoc
// @Summary     List posts
// @Description Returns filtered list of posts for the authenticated user
// @Tags        posts
// @Produce     json
// @Security    BearerAuth
// @Param       name    query    string false "Filter by name"
// @Param       content query    string false "Filter by content"
// @Success     200 {array}		 responses.FindOneResponse
// @Failure     401 {object}     responses.ErrorResponse "Unauthorized"
// @Failure     500 {object}     responses.ErrorResponse "Internal server error"
// @Router      /posts [get]
func (h *FindListPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserFromRequest(r)
	if user == nil {
		httputils.ResponseError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	findListDTO := &dto.FindListDTO{
		Name:    r.URL.Query().Get("name"),
		Content: r.URL.Query().Get("content"),
		UserId:  user.Id,
	}

	postList, err := h.postService.FindList(findListDTO)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.ResponseOk(w, http.StatusOK, responses.NewFindListResponse(postList))
}
