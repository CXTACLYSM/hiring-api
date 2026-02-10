package handlers

import (
	"errors"
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/auth/shared/infrastructure/httputils"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/services"
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
	claims := httputils.ClaimsFromRequest(r)
	if claims == nil {
		httputils.WriteError(w, errors.New("claims is nil"))
		return
	}

	user, err := h.userService.Me(claims.UserId)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}
	if user == nil {
		httputils.ResponseError(w, http.StatusNotFound, "404 not found.")
		return
	}

	httputils.ResponseSuccess(w, http.StatusOK, map[string]any{
		"id":       user.Id,
		"username": user.Username,
		"email":    user.Email,
	})
}
