package handlers

import (
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/services"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/infrastructure/responses"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/httputils"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/middlewares"
)

type MeHandler struct {
	userService *services.UserService
}

func NewMeHandler(userService *services.UserService) *MeHandler {
	return &MeHandler{
		userService: userService,
	}
}

// ServeHTTP godoc
// @Summary     Get current user
// @Description Returns authenticated user profile
// @Tags        user
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} responses.MeResponse
// @Failure     401 {object} responses.ErrorResponse "Unauthorized"
// @Failure     404 {object} responses.ErrorResponse "User not found"
// @Failure     500 {object} responses.ErrorResponse "Internal server error"
// @Router      /me [get]
func (h *MeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserFromRequest(r)
	if user == nil {
		httputils.ResponseError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	freshUser, err := h.userService.Me(user.Id)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}
	if freshUser == nil {
		httputils.ResponseError(w, http.StatusNotFound, "not found")
		return
	}

	httputils.ResponseOk(w, http.StatusOK, responses.NewMeResponse(freshUser))
}
