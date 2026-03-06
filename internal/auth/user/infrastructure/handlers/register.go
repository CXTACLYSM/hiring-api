package handlers

import (
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/dto"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/services"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/infrastructure/responses"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/httputils"
)

type RegisterHandler struct {
	authService *services.AuthService
}

func NewRegisterHandler(authService *services.AuthService) *RegisterHandler {
	return &RegisterHandler{
		authService: authService,
	}
}

// ServeHTTP godoc
// @Summary     Register new user
// @Description Creates a new user account and returns JWT token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body     dto.RegisterDTO true "Registration data"
// @Success     201 {object}     responses.RegisterResponse
// @Failure     422 {object}     responses.ValidationErrorsResponse "Validation errors"
// @Failure     422 {object}     responses.ErrorResponse            "User already exists"
// @Failure     500 {object}     responses.ErrorResponse            "Internal server error"
// @Router      /register [post]
func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var registerDTO dto.RegisterDTO
	if err := httputils.DecodeJSON(r, &registerDTO); err != nil {
		httputils.WriteError(w, err)
		return
	}

	token, err := h.authService.Register(r.Context(), &registerDTO)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.ResponseOk(w, http.StatusCreated, responses.NewRegisterResponse(token))
}
