package handlers

import (
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/dto"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/services"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/infrastructure/responses"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/httputils"
)

type LoginHandler struct {
	authService *services.AuthService
}

func NewLoginHandler(authService *services.AuthService) *LoginHandler {
	return &LoginHandler{
		authService: authService,
	}
}

// ServeHTTP godoc
// @Summary     Login
// @Description Authenticates user by username/email and password, returns JWT token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body     dto.LoginDTO true "Login credentials"
// @Success     200 {object}     responses.LoginResponse
// @Failure     422 {object}     responses.ErrorResponse "Incorrect login or password"
// @Failure     500 {object}     responses.ErrorResponse "Internal server error"
// @Router      /login [post]
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var loginDTO dto.LoginDTO
	if err := httputils.DecodeJSON(r, &loginDTO); err != nil {
		httputils.WriteError(w, err)
		return
	}

	token, err := h.authService.Login(&loginDTO)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.ResponseOk(w, http.StatusOK, responses.NewLoginResponse(token))
}
