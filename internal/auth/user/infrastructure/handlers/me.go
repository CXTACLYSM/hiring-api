package handlers

import (
	"errors"
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

func (h *MeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserFromRequest(r)
	if user == nil {
		httputils.WriteError(w, errors.New("claims is nil"))
		return
	}

	freshUser, err := h.userService.Me(user.Id)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}
	if freshUser == nil {
		httputils.ResponseError(w, http.StatusNotFound, "404 not found.")
		return
	}

	httputils.ResponseOk(w, http.StatusOK, responses.NewMeResponse(freshUser))
}
