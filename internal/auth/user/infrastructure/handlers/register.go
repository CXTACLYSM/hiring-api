package handlers

import (
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/auth/shared/infrastructure/httputils"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/dto"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/services"
)

type RegisterHandler struct {
	authService *services.AuthService
}

func NewRegisterHandler(authService *services.AuthService) *RegisterHandler {
	return &RegisterHandler{
		authService: authService,
	}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var registerDTO *dto.RegisterDTO
	if err := httputils.DecodeJSON(r, &registerDTO); err != nil {
		httputils.WriteError(w, err)
		return
	}

	token, err := h.authService.Register(registerDTO)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.ResponseSuccess(w, http.StatusCreated, map[string]string{"token": token})
}
