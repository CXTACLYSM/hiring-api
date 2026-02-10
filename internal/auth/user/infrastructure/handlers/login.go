package handlers

import (
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/auth/shared/infrastructure/httputils"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/dto"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/services"
)

type LoginHandler struct {
	authService *services.AuthService
}

func NewLoginHandler(authService *services.AuthService) *LoginHandler {
	return &LoginHandler{
		authService: authService,
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var loginDTO *dto.LoginDTO
	if err := httputils.DecodeJSON(r, &loginDTO); err != nil {
		httputils.WriteError(w, err)
		return
	}

	token, err := h.authService.Login(loginDTO)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.ResponseSuccess(w, http.StatusOK, map[string]string{"token": token})
}
